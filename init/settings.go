package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerSetting `mapstructure:"server"`
	Mongo  MongoSetting  `mapstructure:"mongo"`
	Redis  RedisSetting  `mapstructure:"redis"`
	//Grafana   GrafanaSetting  `mapstructure:"grafana"`
	//Logger    LogSetting      `mapstructure:"logger"`
	//Resend    ResendSetting   `mapstructure:"resend"`
	// Kafka KafkaSetting `mapstructure:"kafka"`
}

type ServerSetting struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}
type MongoSetting struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}
type RedisSetting struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

func LoadConfig() Config {
	return LoadConfigFrom(".")
}

func LoadConfigFrom(path string) Config {
	viper := viper.New()
	viper.AddConfigPath(path)

	// also search next to the executable (useful when binary is run from a different cwd)
	if exe, err := os.Executable(); err == nil {
		viper.AddConfigPath(filepath.Dir(exe))
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Failed to read the configuration %w \n", err))
	}
	fmt.Println("Server Port:: ", viper.GetInt("server.port"))
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Unable to decode configuration %v", err)
	}

	return config
}
