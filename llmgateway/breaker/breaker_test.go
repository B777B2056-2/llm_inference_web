package breaker

import (
	"errors"
	"llm_online_interence/llmgateway/confparser"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// 测试配置
var testConfig = confparser.BreakerConf{
	Enable:                   true,
	MaxFailures:              3,
	OpenStateTimeInSeconds:   1,
	HalfOpenStateMaxRequests: 2,
	HalfOpenSuccessThreshold: 2,
}

func TestClosedState(t *testing.T) {
	b := newStateMachine(testConfig)

	// 初始状态应为Closed
	if state := b.currentState(); state.name() != closed {
		t.Fatalf("初始状态错误，期望 Closed，实际 %s", state)
	}
	// 失败次数未达阈值
	for i := 0; i < int(testConfig.MaxFailures)-1; i++ {
		if !b.allowRequest() {
			t.Error("应允许请求")
		}
		b.updateWhenRequestFailed()
	}
	if state := b.currentState(); state.name() != closed {
		t.Fatal("不应触发状态转换")
	}

	// 触发失败阈值
	b.updateWhenRequestFailed()
	if state := b.currentState(); state.name() != open {
		t.Fatal("应转换到Open状态")
	}
}

func TestOpenState(t *testing.T) {
	b := newStateMachine(testConfig)
	b.transition(open)

	// Open状态应拒绝请求
	if b.allowRequest() {
		t.Error("Open状态应拒绝请求")
	}

	// 等待超时时间
	time.Sleep(time.Duration(testConfig.OpenStateTimeInSeconds+1) * time.Second)
	if !b.allowRequest() {
		t.Error("超时后应允许请求")
	}
	if state := b.currentState(); state.name() != halfOpen {
		t.Fatal("应转换到HalfOpen状态")
	}
}

func TestHalfOpenState(t *testing.T) {
	b := newStateMachine(testConfig)
	b.transition(halfOpen)

	// 测试成功场景
	for i := 0; i < int(testConfig.HalfOpenStateMaxRequests); i++ {
		if !b.allowRequest() {
			t.Errorf("应允许请求，但第%d次请求被拒绝", i+1)
		}
		b.updateWhenRequestSuccess()
	}
	if state := b.currentState(); state.name() != closed {
		t.Fatal("成功阈值达成应转换到Closed")
	}

	// 测试失败场景
	b.transition(halfOpen)
	for i := 0; i < int(testConfig.HalfOpenStateMaxRequests)-1; i++ {
		b.updateWhenRequestFailed()
	}
	if state := b.currentState(); state.name() != open {
		t.Fatal("半开状态失败应转换到Open")
	}
}

func TestConcurrentAccess(t *testing.T) {
	b := newStateMachine(testConfig)
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
	if b.currentState().name() != open {
		t.Error("并发访问后应处于Open状态")
	}
}
