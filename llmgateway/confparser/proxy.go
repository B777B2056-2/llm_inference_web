package confparser

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type BreakerConf struct {
	Enable                   bool   `yaml:"enable"`
	MaxFailures              uint64 `yaml:"maxFailures"`
	OpenStateTimeInSeconds   uint64 `yaml:"openStateTimeInSeconds"`
	HalfOpenStateMaxRequests uint64 `yaml:"halfOpenStateMaxRequests"`
	HalfOpenSuccessThreshold uint64 `yaml:"halfOpenSuccessThreshold"`
}

type RateLimitURLConfItem struct {
	URI            string `yaml:"uri"`
	BucketSize     uint   `yaml:"bucketSize"`
	TokenPerSecond uint   `yaml:"tokenPerSecond"`
}

type BackendConfigItem struct {
	SvcName               string                 `yaml:"svcName"`
	GroupName             string                 `yaml:"groupName"`
	NeedRefreshToken      bool                   `yaml:"needRefreshToken"`
	Protocol              string                 `yaml:"protocol"`
	LoadBalanceStrategy   string                 `yaml:"loadBalanceStrategy"`
	ConnectionTimeout     int                    `yaml:"connectionTimeout"`
	ResponseTimeout       int                    `yaml:"responseTimeout"`
	Breaker               BreakerConf            `yaml:"breaker"`
	NeedAuthURLs          []string               `yaml:"needAuthURLs"`
	NeedRateLimitURLConf  []RateLimitURLConfItem `yaml:"needRateLimitURLConf"`  // url层面限流配置（限制总体并发）
	NeedRateLimitUserConf []RateLimitURLConfItem `yaml:"needRateLimitUserConf"` // 用户层面限流配置（限制单用户并发）
}

type GRPCConfigItem struct {
	Host string `yaml:"host"`
	Port uint   `yaml:"port"`
}

var ProxyConfig struct {
	BlackedIPs      []string                                   `yaml:"blackedIPs"`
	Backends        []BackendConfigItem                        `yaml:"backends"`
	GRPC            map[string]GRPCConfigItem                  `yaml:"grpc"`
	SvcURIRateLimit map[string]map[string]RateLimitURLConfItem `yaml:"-"` // 外层key为svc name，内层key为uri
	UserRateLimit   map[string]map[string]RateLimitURLConfItem `yaml:"-"` // 外层key为svc name，内层key为uri
}

func InitProxyConfig(filePath string) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Errorf("read config file [%s] error: %v", filePath, err))
	}
	err = yaml.Unmarshal(file, &ProxyConfig)
	if err != nil {
		panic(fmt.Errorf("unmarshal config file [%s] error: %v", filePath, err))
	}

	ProxyConfig.SvcURIRateLimit = make(map[string]map[string]RateLimitURLConfItem)
	ProxyConfig.UserRateLimit = make(map[string]map[string]RateLimitURLConfItem)
	for idx := range ProxyConfig.Backends {
		ProxyConfig.Backends[idx].Protocol = strings.ToLower(ProxyConfig.Backends[idx].Protocol)

		ProxyConfig.SvcURIRateLimit[ProxyConfig.Backends[idx].SvcName] = make(map[string]RateLimitURLConfItem)
		ProxyConfig.UserRateLimit[ProxyConfig.Backends[idx].SvcName] = make(map[string]RateLimitURLConfItem)
		for _, urlConf := range ProxyConfig.Backends[idx].NeedRateLimitURLConf {
			ProxyConfig.SvcURIRateLimit[ProxyConfig.Backends[idx].SvcName][urlConf.URI] = urlConf
		}
		for _, urlConf := range ProxyConfig.Backends[idx].NeedRateLimitUserConf {
			ProxyConfig.UserRateLimit[ProxyConfig.Backends[idx].SvcName][urlConf.URI] = urlConf
		}
	}
}
