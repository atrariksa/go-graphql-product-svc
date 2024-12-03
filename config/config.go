package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig    ServerConfig    `mapstructure:"server"`
	RedisConfig     RedisConfig     `mapstructure:"redis"`
	DBConfig        DBConfig        `mapstructure:"database"`
	AuthTokenConfig AuthTokenConfig `mapstructure:"auth-token"`
}

type ServerConfig struct {
	Host string 	`mapstructure:"host"`
	Port int    	`mapstructure:"port"`
	Env  string		`mapstructure:"env"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type DBConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
}

type AuthTokenConfig struct {
	Duration  time.Duration `mapstructure:"duration"`
	SecretKey string        `mapstructure:"secretkey"`
}

func GetConfig() *Config {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config.yaml")
	v.AddConfigPath("./config")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = v.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

