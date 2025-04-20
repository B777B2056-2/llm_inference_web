package resource

func Init() {
	initLogger()
	initRedis()
	initMySQL()
	initValidator()
}
