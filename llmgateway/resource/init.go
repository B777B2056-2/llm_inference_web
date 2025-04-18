package resource

import "sync"

var once sync.Once

func Init() {
	once.Do(func() {
		initLogger()
		initRedis()
	})
}
