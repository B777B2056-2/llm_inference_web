package breaker

import (
	"llm_online_interence/llmgateway/confparser"
	"sync"
)

var (
	once   sync.Once
	fsmMap map[string]*stateMachine
)

func Init() {
	once.Do(func() {
		fsmMap = make(map[string]*stateMachine)
		for _, backend := range confparser.ProxyConfig.Backends {
			if !backend.Breaker.Enable {
				continue
			}
			fsmMap[backend.SvcName] = newStateMachine(backend.Breaker)
		}
	})
}

func GetBreakerBySvcName(svcName string) *stateMachine {
	return fsmMap[svcName]
}
