package bootstrap

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerSetting `mapstructure:"server"`
	Mongo  MongoSetting  `mapstructure:"mongo"`
	Redis  RedisSetting  `mapstructure:"redis"`
	//Grafana   GrafanaSetting  `mapstructure:"grafana"`
	//Logger    LogSetting      `mapstructure:"logger"`
	//Resend    ResendSetting   `mapstructure:"resend"`
	Kafka KafkaSetting `mapstructure:"kafka"`
}

type ServerSetting struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}
type MongoSetting struct {
	URI  string `mapstructure:"URI"`
	Port string `mapstructure:"port"`
}
type RedisSetting struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}
type KafkaSetting struct {
}

func LoadConfig() Config {
	viper := viper.New()
	viper.AddConfigPath("./config")

	viper.SetConfigName("local")
	viper.SetConfigFile("yml")

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
