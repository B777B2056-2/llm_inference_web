package breaker

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/sirupsen/logrus"
	"llm_online_interence/llmgateway/confparser"
	"llm_online_interence/llmgateway/resource"
	"time"
)

const redisKeyPrefix = "breaker_"

// 断路器名称定义
const (
	closed   = "closed"
	open     = "open"
	halfOpen = "halfOpen"
)

// 断路器状态定义
type abstractState interface {
	allowRequest(tx *redis.Tx) (bool, error)
	name() string
	enter(tx *redis.Tx) error
	doAction(tx *redis.Tx) error
}

// 关闭状态，允许请求通过，同时记录失败次数
type closedState struct {
	ctx         context.Context
	redisKey    string
	redisField  string
	maxFailures uint64
}

func newClosedState(ctx context.Context, redisKey string, maxFailures uint64) *closedState {
	c := &closedState{ctx: ctx, redisKey: redisKey, maxFailures: maxFailures}
	redisField := c.name() + "#" + "FailuresCounter"
	// 仅当键不存在时初始化
	if err := resource.RedisClient.HSetNX(ctx, redisKey, redisField, 0).Err(); err != nil {
		panic(errors.New(redisKey + ": " + err.Error()))
	}
	c.redisField = redisField
	return c
}

func (c *closedState) allowRequest(tx *redis.Tx) (bool, error) {
	failuresCounter, err := tx.HGet(c.ctx, c.redisKey, c.redisField).Uint64()
	if err != nil {
		return false, err
	}
	if failuresCounter >= c.maxFailures {
		return false, nil
	}
	return true, nil
}

func (c *closedState) name() string { return closed }

func (c *closedState) enter(tx *redis.Tx) error {
	return tx.HSet(c.ctx, c.redisKey, c.redisField, 0).Err()
}

func (c *closedState) doAction(tx *redis.Tx) error {
	return tx.HIncrBy(c.ctx, c.redisKey, c.redisField, 1).Err()
}

// 开启状态，在指定时间内，所有请求被拒绝
type openState struct {
	ctx                 context.Context
	redisKey            string
	redisField          string
	stateDurationSecond uint64
}

func newOpenState(ctx context.Context, redisKey string, stateDurationSecond uint64) *openState {
	o := &openState{ctx: ctx, redisKey: redisKey, stateDurationSecond: stateDurationSecond}
	redisField := o.name() + "#" + "StartTimestamp"
	// 仅当键不存在时初始化
	if err := resource.RedisClient.HSetNX(ctx, redisKey, redisField, 0).Err(); err != nil {
		panic(errors.New(redisKey + ": " + err.Error()))
	}
	o.redisField = redisField
	return o
}

func (o *openState) allowRequest(tx *redis.Tx) (bool, error) {
	// 在断路器打开状态持续时间内，不允许请求通过
	redisTime, err := tx.Time(o.ctx).Result()
	if err != nil {
		return false, err
	}
	now := redisTime.Unix()
	startTime, err := tx.HGet(o.ctx, o.redisKey, o.redisField).Int64()
	if err != nil {
		return false, err
	}
	if uint64(now-startTime) <= o.stateDurationSecond {
		return false, nil
	}
	return true, nil
}

func (o *openState) name() string { return open }

// 设置所有新请求的超时时间
func (o *openState) enter(tx *redis.Tx) error {
	redisTime, err := tx.Time(o.ctx).Result()
	if err != nil {
		return err
	}
	return tx.HSet(o.ctx, o.redisKey, o.redisField, redisTime.Unix()).Err()
}

func (o *openState) doAction(_ *redis.Tx) error { return nil }

// 半开状态，允许n个请求通过，并记录成功次数，其他请求被拒绝
type halfOpenState struct {
	ctx                       context.Context
	redisKey                  string
	totalReqCounterRedisField string
	successCounterRedisField  string
	maxReqThreshold           uint64
	successThreshold          uint64
}

func newHalfOpenState(ctx context.Context, redisKey string,
	maxReqThreshold, successThreshold uint64) *halfOpenState {
	h := &halfOpenState{
		ctx:              ctx,
		redisKey:         redisKey,
		maxReqThreshold:  maxReqThreshold,
		successThreshold: successThreshold,
	}
	successCounterRedisField := h.name() + "#" + "SuccessCounter"
	totalReqCounterRedisField := h.name() + "#" + "TotalReqCounter"
	// 仅当键不存在时初始化
	if err := resource.RedisClient.HSetNX(ctx, redisKey, successCounterRedisField, 0).Err(); err != nil {
		panic(errors.New(redisKey + ": " + err.Error()))
	}
	if err := resource.RedisClient.HSetNX(ctx, redisKey, totalReqCounterRedisField, 0).Err(); err != nil {
		panic(errors.New(redisKey + ": " + err.Error()))
	}
	h.totalReqCounterRedisField = totalReqCounterRedisField
	h.successCounterRedisField = successCounterRedisField
	return h
}

func (h *halfOpenState) allowRequest(tx *redis.Tx) (bool, error) {
	totalReqCounter, err := tx.HGet(h.ctx, h.redisKey, h.totalReqCounterRedisField).Uint64()
	if err != nil {
		return false, nil
	}
	if totalReqCounter < h.maxReqThreshold {
		return true, nil
	}
	return false, nil
}

func (h *halfOpenState) name() string { return halfOpen }

func (h *halfOpenState) enter(tx *redis.Tx) error {
	return tx.HMSet(
		h.ctx, h.redisKey,
		h.totalReqCounterRedisField, 0,
		h.successCounterRedisField, 0,
	).Err()
}

func (h *halfOpenState) doAction(tx *redis.Tx) error {
	return tx.HIncrBy(h.ctx, h.redisKey, h.totalReqCounterRedisField, 1).Err()
}

func (h *halfOpenState) incrSuccessCounter(tx *redis.Tx) error {
	return tx.HIncrBy(h.ctx, h.redisKey, h.successCounterRedisField, 1).Err()
}

func (h *halfOpenState) canTransition2Closed(tx *redis.Tx) (bool, error) {
	successCounter, err := tx.HGet(h.ctx, h.redisKey, h.successCounterRedisField).Uint64()
	if err != nil {
		return false, err
	}
	if successCounter < h.successThreshold {
		return false, nil
	}
	return true, nil
}

// 状态机，控制状态流转
type stateMachine struct {
	ctx                   context.Context
	redisMutex            *redsync.Mutex
	redisKey              string
	curStateIdxRedisField string
	stateList             []abstractState
}

func newStateMachine(svcName string, conf confparser.BreakerConf) *stateMachine {
	ctx := context.Background()
	redisKey := fmt.Sprintf("%s%s", redisKeyPrefix, svcName)
	s := &stateMachine{
		ctx: ctx,
		redisMutex: resource.RedisLocker.NewMutex(
			svcName,
			redsync.WithExpiry(24*time.Hour),
			redsync.WithRetryDelay(
				time.Duration(confparser.ResourceConfig.Redis.Lock.RetryDelayMs)*time.Millisecond,
			),
			redsync.WithTries(confparser.ResourceConfig.Redis.Lock.MaxRetries),
		),
		stateList: []abstractState{
			newClosedState(
				ctx,
				redisKey,
				conf.MaxFailures,
			),
			newOpenState(
				ctx,
				redisKey,
				conf.OpenStateTimeInSeconds,
			),
			newHalfOpenState(
				ctx,
				redisKey,
				conf.HalfOpenStateMaxRequests,
				conf.HalfOpenSuccessThreshold,
			),
		},
		redisKey:              redisKey,
		curStateIdxRedisField: "state_machine_current_state_idx",
	}
	if err := resource.RedisClient.HSetNX(ctx, redisKey, s.curStateIdxRedisField, 0).Err(); err != nil {
		panic(errors.New(redisKey + ": " + err.Error()))
	}
	return s
}

func (s *stateMachine) transition(tx *redis.Tx, stateName string) (int, error) {
	currentIdx := s.currentStateIdx(tx)
	currentState := s.stateList[currentIdx]
	if currentState.name() == stateName {
		return currentIdx, nil
	}
	for i, state := range s.stateList {
		if state.name() == stateName {
			if err := tx.HSet(s.ctx, s.redisKey, s.curStateIdxRedisField, i).Err(); err != nil {
				return currentIdx, err
			}
			if err := s.stateList[i].enter(tx); err != nil {
				return currentIdx, err
			}
			return i, nil
		}
	}
	return -1, nil
}

func (s *stateMachine) currentStateIdx(tx *redis.Tx) int {
	curStateIdx, err := tx.HGet(s.ctx, s.redisKey, s.curStateIdxRedisField).Int()
	if err != nil {
		curStateIdx = 1 // redis出错时的降级策略：返回open状态，拒绝所有请求
	}
	return curStateIdx
}

// allowRequest 检查是否允许请求通过，如果当前状态是打开状态，则切换到半开状态
func (s *stateMachine) allowRequest(tx *redis.Tx) (bool, error) {
	currentIdx := s.currentStateIdx(tx)
	newStateIdx := currentIdx
	current := s.stateList[currentIdx]
	if current.name() == open {
		allowed, err := current.allowRequest(tx)
		if err != nil {
			return false, err
		}
		if allowed {
			newStateIdx, err = s.transition(tx, halfOpen)
			if err != nil {
				return false, err
			}
		}
	}
	return s.stateList[newStateIdx].allowRequest(tx)
}

// updateWhenRequestFailed 请求失败，计数器加1，如果计数器达到阈值，则切换到打开状态
func (s *stateMachine) updateWhenRequestFailed(tx *redis.Tx) error {
	switch currentState := s.stateList[s.currentStateIdx(tx)]; currentState.name() {
	// 关闭状态
	case closed:
		// 请求失败，计数器加1
		if err := currentState.doAction(tx); err != nil {
			return err
		}
		// 请求失败，计数器达到阈值，切换到打开状态
		allow, err := currentState.allowRequest(tx)
		if err != nil {
			return err
		}
		if !allow {
			if _, err := s.transition(tx, open); err != nil {
				return err
			}
		}
	case halfOpen:
		if err := currentState.doAction(tx); err != nil {
			return err
		}
		if _, err := s.transition(tx, open); err != nil {
			return err
		}
	}
	return nil
}

// updateWhenRequestSuccess 请求成功，计数器加1，如果计数器达到阈值，则切换到关闭状态
func (s *stateMachine) updateWhenRequestSuccess(tx *redis.Tx) error {
	// 非半开状态，不处理
	currentState := s.stateList[s.currentStateIdx(tx)]
	if currentState.name() != halfOpen {
		return nil
	}
	// 半开状态，总请求计数器加1
	if err := currentState.doAction(tx); err != nil {
		return err
	}
	// 半开状态，成功请求计数器加1
	state := currentState.(*halfOpenState)
	if err := state.incrSuccessCounter(tx); err != nil {
		return err
	}
	// 请求成功数大于设定的阈值时，切换到关闭状态
	canPass, err := state.canTransition2Closed(tx)
	if err != nil {
		return err
	}
	if canPass {
		if _, err := s.transition(tx, closed); err != nil {
			return err
		}
	}
	return nil
}

func (s *stateMachine) Execute(ctx *gin.Context) bool {
	if err := s.redisMutex.Lock(); err != nil {
		return false
	}
	defer func() {
		_, err := s.redisMutex.Unlock()
		if err != nil {
			resource.Logger.WithFields(
				logrus.Fields{
					"error": err,
				},
			).Error("failed to unlock redis")
		}
	}()

	allowRequest := false
	err := resource.RedisClient.Watch(ctx, func(tx *redis.Tx) (err error) {
		if hasErr := len(ctx.Errors) != 0; hasErr {
			if err := s.updateWhenRequestFailed(tx); err != nil {
				return err
			}
		} else {
			if err := s.updateWhenRequestSuccess(tx); err != nil {
				return err
			}
		}
		allowRequest, err = s.allowRequest(tx)
		return
	}, s.redisKey)
	if err != nil {
		allowRequest = false
		resource.Logger.WithFields(
			logrus.Fields{
				"redisKey": s.redisKey,
				"error":    err,
			},
		).Error("failed to execute breaker state machine")
	}
	return allowRequest
}
