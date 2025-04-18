package breaker

import (
	"llm_online_interence/llmgateway/confparser"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 断路器名称定义
const (
	closed   = "closed"
	open     = "open"
	halfOpen = "halfOpen"
)

// 断路器状态定义
type abstractState interface {
	allowRequest() bool
	name() string
	enter()
	doAction()
	leave()
}

// 关闭状态，允许请求通过，同时记录失败次数
type closedState struct {
	failedCounter uint64
	maxFailures   uint64
}

func newClosedState(maxFailures uint64) *closedState {
	return &closedState{maxFailures: maxFailures}
}

func (c *closedState) allowRequest() bool {
	if c.failedCounter < c.maxFailures {
		return true
	}
	return false
}

func (c *closedState) name() string { return closed }
func (c *closedState) enter()       { c.failedCounter = 0 }
func (c *closedState) doAction()    { c.failedCounter++ }
func (c *closedState) leave()       {}

// 开启状态，在指定时间内，所有请求被拒绝
type openState struct {
	startTime           uint64
	stateDurationSecond uint64
}

func newOpenState(stateDurationSecond uint64) *openState {
	return &openState{stateDurationSecond: stateDurationSecond}
}

func (o *openState) allowRequest() bool {
	// 在断路器打开状态持续时间内，不允许请求通过
	if uint64(time.Now().Unix())-o.startTime <= o.stateDurationSecond {
		return false
	}
	return true
}

func (o *openState) name() string { return open }

// 设置所有新请求的超时时间
func (o *openState) enter() {
	o.startTime = uint64(time.Now().Unix())
}

func (o *openState) doAction() {}

// 重置所有请求的超时时间
func (o *openState) leave() { o.stateDurationSecond = 0 }

// 半开状态，允许n个请求通过，并记录成功次数，其他请求被拒绝
type halfOpenState struct {
	totalReqCounter  uint64
	maxReqThreshold  uint64
	successCounter   uint64
	successThreshold uint64
}

func newHalfOpenState(maxReqThreshold, successThreshold uint64) *halfOpenState {
	return &halfOpenState{
		maxReqThreshold:  maxReqThreshold,
		successThreshold: successThreshold,
	}
}

func (h *halfOpenState) allowRequest() bool         { return h.totalReqCounter < h.maxReqThreshold }
func (h *halfOpenState) name() string               { return halfOpen }
func (h *halfOpenState) enter()                     { h.successCounter = 0; h.totalReqCounter = 0 }
func (h *halfOpenState) doAction()                  { h.totalReqCounter++ }
func (h *halfOpenState) leave()                     { h.successCounter = 0; h.totalReqCounter = 0 }
func (h *halfOpenState) incrSuccessCounter()        { h.successCounter++ }
func (h *halfOpenState) canTransition2Closed() bool { return h.successCounter >= h.successThreshold }

// 状态机，控制状态流转
type stateMachine struct {
	mu          sync.Mutex
	curStateIdx int
	stateList   []abstractState
}

func newStateMachine(conf confparser.BreakerConf) *stateMachine {
	return &stateMachine{
		stateList: []abstractState{
			newClosedState(conf.MaxFailures),
			newOpenState(conf.OpenStateTimeInSeconds),
			newHalfOpenState(
				conf.HalfOpenStateMaxRequests,
				conf.HalfOpenSuccessThreshold,
			),
		},
		curStateIdx: 0,
	}
}

func (s *stateMachine) transition(stateName string) {
	for i, state := range s.stateList {
		if state.name() == stateName {
			s.currentState().leave()
			s.curStateIdx = i
			s.currentState().enter()
			return
		}
	}
}

func (s *stateMachine) currentState() abstractState {
	return s.stateList[s.curStateIdx]
}

// allowRequest 检查是否允许请求通过，如果当前状态是打开状态，则切换到半开状态
func (s *stateMachine) allowRequest() bool {
	current := s.currentState()
	if current.name() == open && current.allowRequest() {
		s.transition(halfOpen)
	}
	return s.currentState().allowRequest()
}

// updateWhenRequestFailed 请求失败，计数器加1，如果计数器达到阈值，则切换到打开状态
func (s *stateMachine) updateWhenRequestFailed() {
	switch s.currentState().name() {
	// 关闭状态
	case closed:
		// 请求失败，计数器加1
		s.currentState().doAction()
		// 请求失败，计数器达到阈值，切换到打开状态
		if !s.currentState().allowRequest() {
			s.transition(open)
		}
	case halfOpen:
		s.currentState().doAction()
		s.transition(open)
	}
}

// updateWhenRequestSuccess 请求成功，计数器加1，如果计数器达到阈值，则切换到关闭状态
func (s *stateMachine) updateWhenRequestSuccess() {
	// 非半开状态，不处理
	if s.currentState().name() != halfOpen {
		return
	}
	// 半开状态，总请求计数器加1
	s.currentState().doAction()
	// 半开状态，成功请求计数器加1
	state := s.currentState().(*halfOpenState)
	state.incrSuccessCounter()
	// 请求成功数大于设定的阈值时，切换到关闭状态
	if state.canTransition2Closed() {
		s.transition(closed)
	}
}

func (s *stateMachine) Execute(ctx *gin.Context) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if hasErr := len(ctx.Errors) != 0; hasErr {
		s.updateWhenRequestFailed()
	} else {
		s.updateWhenRequestSuccess()
	}
	return s.allowRequest()
}
