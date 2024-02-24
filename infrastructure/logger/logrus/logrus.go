package logrus

import (
	"fmt"
	"jwt-auth-microservice/config"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	appName    string
	appVersion string
)

func InitializeLogrusLogger(conf config.EnvironmentConfig) {
	appName = conf.App.Name
	appVersion = conf.App.Version

	currentTime := time.Now()
	date := currentTime.Format("20060102")
	path := fmt.Sprintf("%s/%s-%s.%s", conf.Log.Path, conf.Log.Prefix, date, conf.Log.Ext)

	log.SetFormatter(&log.JSONFormatter{})

	err := os.MkdirAll(filepath.Dir(path), 0770)
	if err == nil {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.SetOutput(os.Stdout)
			return
		}
		log.SetOutput(file)
	} else {
		log.SetOutput(os.Stdout)
	}
}

func LogInfo(logtype, message string) {
	log.WithFields(log.Fields{
		"app_name":    appName,
		"app_version": appVersion,
		"log_type":    logtype,
	}).Info(message)
}

func LogInfoWithData(data interface{}, logtype, message string) {
	log.WithFields(log.Fields{
		"app_name":    appName,
		"app_version": appVersion,
		"data":        data,
		"log_type":    logtype,
	}).Info(message)
}

func LogError(logtype, errtype, message string) {
	log.WithFields(log.Fields{
		"app_name":    appName,
		"app_version": appVersion,
		"error_type":  errtype,
		"log_type":    logtype,
	}).Error(message)
}
