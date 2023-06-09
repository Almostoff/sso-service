package loggerService

import "AuthService/config"

var instance Logger = nil

func Init(cfg *config.Config) {
	instance = NewLoggerClient(cfg)
}

func GetInstance() Logger {
	return instance
}
