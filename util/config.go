package util

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	DbSource             string        `mapstructure:"DB_SOURCE"`
	ServerAddr           string        `mapstructure:"SERVER_ADDR"`
	TokenSymmetryKey     string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	GrpcServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	Environment          string        `mapstructure:"ENVIRONMENT"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
