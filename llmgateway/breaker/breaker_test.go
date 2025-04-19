package breaker

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"llm_online_interence/llmgateway/confparser"
	"llm_online_interence/llmgateway/resource"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// 测试配置
const testSvcName = "testSvc"

var testConfig = confparser.BreakerConf{
	Enable:                   true,
	MaxFailures:              3,
	OpenStateTimeInSeconds:   1,
	HalfOpenStateMaxRequests: 2,
	HalfOpenSuccessThreshold: 2,
}

func currentState(b *stateMachine, tx *redis.Tx) abstractState {
	return b.stateList[b.currentStateIdx(tx)]
}

func initResource() {
	confparser.ResourceConfig.Logger.Level = resource.LoggerLevelDebug
	confparser.ResourceConfig.Logger.OutPutPath = "./logs"
	confparser.ResourceConfig.Logger.MaxSingleFileSizeMB = 100
	confparser.ResourceConfig.Logger.MaxBackups = 1
	confparser.ResourceConfig.Logger.MaxStorageAgeInDays = 1

	confparser.ResourceConfig.Redis.Host = "127.0.0.1"
	confparser.ResourceConfig.Redis.Port = 6379
	confparser.ResourceConfig.Redis.Password = "123456"
	confparser.ResourceConfig.Redis.DB = 0
	confparser.ResourceConfig.Redis.PoolSize = 4
	confparser.ResourceConfig.Redis.DialTimeoutSecond = 5
	confparser.ResourceConfig.Redis.ReadTimeoutSecond = 5
	confparser.ResourceConfig.Redis.WriteTimeoutSecond = 5
	confparser.ResourceConfig.Redis.ConnMaxRetries = 3
	confparser.ResourceConfig.Redis.TxMaxRetries = 2
	confparser.ResourceConfig.Redis.Lock.RetryDelayMs = 250
	confparser.ResourceConfig.Redis.Lock.MaxRetries = 2
	resource.Init()
}

func TestClosedState(t *testing.T) {
	initResource()
	b := newStateMachine(testSvcName, testConfig)

	if err := resource.RedisClient.Watch(b.ctx, func(tx *redis.Tx) (err error) {
		if _, err := b.transition(tx, closed); err != nil {
			return err
		}
		// 初始状态应为Closed
		if state := currentState(b, tx); state.name() != closed {
			return fmt.Errorf("初始状态错误，期望 Closed，实际 %s", state)
		}
		// 失败次数未达阈值
		for i := 0; i < int(testConfig.MaxFailures)-1; i++ {
			allow, err := b.allowRequest(tx)
			if err != nil {
				return err
			}
			if !allow {
				return fmt.Errorf("应允许请求")
			}
			if err := b.updateWhenRequestFailed(tx); err != nil {
				return err
			}
		}
		if state := currentState(b, tx); state.name() != closed {
			return fmt.Errorf("不应触发状态转换")
		}

		// 触发失败阈值
		if err := b.updateWhenRequestFailed(tx); err != nil {
			return err
		}
		if state := currentState(b, tx); state.name() != open {
			return fmt.Errorf("应转换到Open状态")
		}
		return nil
	}, b.redisKey); err != nil {
		t.Fatal(err)
	}
}

func TestOpenState(t *testing.T) {
	initResource()
	b := newStateMachine(testSvcName, testConfig)

	if err := resource.RedisClient.Watch(b.ctx, func(tx *redis.Tx) (err error) {
		if _, err := b.transition(tx, open); err != nil {
			return err
		}
		// Open状态应拒绝请求
		allow, err := b.allowRequest(tx)
		if err != nil {
			return err
		}
		if allow {
			return fmt.Errorf("open状态应拒绝请求")
		}

		// 等待超时时间
		time.Sleep(time.Duration(testConfig.OpenStateTimeInSeconds+1) * time.Second)
		allow, err = b.allowRequest(tx)
		if err != nil {
			return err
		}
		if !allow {
			t.Error("超时后应允许请求")
		}
		if state := currentState(b, tx); state.name() != halfOpen {
			return fmt.Errorf("应转换到HalfOpen状态")
		}
		return nil
	}, b.redisKey); err != nil {
		t.Fatal(err)
	}
}

func TestHalfOpenState(t *testing.T) {
	initResource()
	b := newStateMachine(testSvcName, testConfig)

	if err := resource.RedisClient.Watch(b.ctx, func(tx *redis.Tx) (err error) {
		if _, err := b.transition(tx, halfOpen); err != nil {
			return err
		}

		// 测试成功场景
		for i := 0; i < int(testConfig.HalfOpenStateMaxRequests); i++ {
			allow, err := b.allowRequest(tx)
			if err != nil {
				return err
			}
			if !allow {
				return fmt.Errorf("应允许请求，但第%d次请求被拒绝", i+1)
			}
			if err := b.updateWhenRequestSuccess(tx); err != nil {
				return err
			}
		}
		if state := currentState(b, tx); state.name() != closed {
			return fmt.Errorf("成功阈值达成应转换到Closed")
		}

		// 测试失败场景
		if _, err := b.transition(tx, halfOpen); err != nil {
			return err
		}
		for i := 0; i < int(testConfig.HalfOpenStateMaxRequests)-1; i++ {
			if err := b.updateWhenRequestFailed(tx); err != nil {
				return err
			}
		}
		if state := currentState(b, tx); state.name() != open {
			return fmt.Errorf("半开状态失败应转换到Open")
		}
		return nil
	}, b.redisKey); err != nil {
		t.Fatal(err)
	}
}

func TestMultiSvcAccess(t *testing.T) {
	initResource()

	var wg sync.WaitGroup
	svcNum := 3
	hasError := atomic.Bool{}

	var machines []*stateMachine
	for i := 0; i < svcNum; i++ {
		b := newStateMachine(testSvcName+fmt.Sprintf("%d", i), testConfig)
		machines = append(machines, b)
	}

	gin.SetMode(gin.TestMode)
	for i := 0; i < svcNum; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			b := machines[idx]
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Errors = append(c.Errors, &gin.Error{Err: errors.New("test error")})
			allowed := b.Execute(c)
			if err := resource.RedisClient.Watch(b.ctx, func(tx *redis.Tx) (err error) {
				if !allowed && currentState(b, tx).name() != open {
					return fmt.Errorf("并发访问后应处于Open状态")
				}
				return nil
			}, b.redisKey); err != nil {
				fmt.Println(err)
				hasError.Store(true)
			}
		}(i)
	}

	wg.Wait()
	if hasError.Load() {
		t.Fatal("error")
	}
}

func TestSingleNodeConcurrentAccess(t *testing.T) {
	initResource()
	b := newStateMachine(testSvcName, testConfig)
	var wg sync.WaitGroup

	gin.SetMode(gin.TestMode)
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Errors = append(c.Errors, &gin.Error{Err: errors.New("test error")})
			_ = b.Execute(c)
		}()
	}

	wg.Wait()

	if err := resource.RedisClient.Watch(b.ctx, func(tx *redis.Tx) (err error) {
		if currentState(b, tx).name() != open {
			return fmt.Errorf("并发访问后应处于Open状态")
		}
		return nil
	}, b.redisKey); err != nil {
		t.Fatal(err)
	}
}

func TestMultiNodeConcurrentAccess(t *testing.T) {
	var wg sync.WaitGroup
	nodeNum := 3
	hasError := atomic.Bool{}

	gin.SetMode(gin.TestMode)
	for i := 0; i < nodeNum; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			initResource()
			b := newStateMachine(testSvcName, testConfig)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Errors = append(c.Errors, &gin.Error{Err: errors.New("test error")})
			allowed := b.Execute(c)
			if err := resource.RedisClient.Watch(b.ctx, func(tx *redis.Tx) (err error) {
				if !allowed && currentState(b, tx).name() != open {
					return fmt.Errorf("并发访问后应处于Open状态")
				}
				return nil
			}, b.redisKey); err != nil {
				fmt.Println(err)
				hasError.Store(true)
			}
		}(i)
	}

	wg.Wait()
	if hasError.Load() {
		t.Fatal("error")
	}
}
