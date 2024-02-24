package logger

import (
	"jwt-auth-microservice/config"
	"jwt-auth-microservice/infrastructure/logger/logrus"
)

func SetupLogger(conf config.EnvironmentConfig) {
	logrus.InitializeLogrusLogger(conf)
}

func LogInfo(logtype, message string) {
	logrus.LogInfo(logtype, message)
}

func LogInfoWithData(data interface{}, logtype, message string) {
	logrus.LogInfoWithData(data, logtype, message)
}

func LogError(logtype, errtype, message string) {
	logrus.LogError(logtype, errtype, message)
}
