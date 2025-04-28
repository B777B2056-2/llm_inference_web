package resource

import "context"

func Init() {
	initLogger()
	initRedis()
	initMySQL()
	initMongoDB(context.Background())
	initValidator()
}
