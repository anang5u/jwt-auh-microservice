package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type EnvironmentConfig struct {
	App App
	Log Log
	//Database database.DatabaseConfig
	//RabbitMq rabbitmq.RabbitmqConfig
}

type App struct {
	Name        string
	Description string
	Version     string
	Port        string
}

type Log struct {
	Path   string
	Prefix string
	Ext    string
}

func SetupConfig(paths ...string) (config EnvironmentConfig, err error) {
	path2file := "config"
	if len(paths) > 0 && len(paths[0]) > 0 {
		path2file = paths[0]
	}

	viper.SetConfigName(path2file)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	//viper.SetEnvPrefix("demo")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("unable to initialize viper: %w", err))
	}
	log.Println("viper config initialized...")

	config = EnvironmentConfig{
		App: App{
			Name:        viper.GetString("AppName"),
			Description: viper.GetString("Description"),
			Version:     viper.GetString("Version"),
			Port:        viper.GetString("Port"),
		},
		Log: Log{
			Path:   viper.GetString("Log.Path"),
			Prefix: viper.GetString("Log.Prefix"),
			Ext:    viper.GetString("Log.Ext"),
		},
		/*
			Database: database.DatabaseConfig{
				Host:     os.Getenv("DB_HOST"),
				Name:     os.Getenv("DB_NAME"),
				Username: os.Getenv("DB_USERNAME"),
				Password: os.Getenv("DB_PASSWORD"),
			},
		*/
	}

	return
}
