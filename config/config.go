package config

import (
	"fmt"
	"os"
)

type Config struct {
	MySQLHost     string `env:"MYSQL_HOST"`
	MySQLPort     string `env:"MYSQL_PORT"`
	MySQLUser     string `env:"MYSQL_USER"`
	MySQLPassword string `env:"MYSQL_PASSWORD"`
	MySQLDBName   string `env:"MYSQL_DB_NAME"`

	RedisHost string `env:"REDIS_HOST"`
	RedisPort string `env:"REDIS_PORT"`
}

func LoadConfig() (Config, error) {
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlDBName := os.Getenv("MYSQL_DB_NAME")

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	if mysqlHost == "" || mysqlPort == "" || mysqlUser == "" || mysqlPassword == "" || mysqlDBName == "" {
		return Config{}, fmt.Errorf("MySQL environment variables are not set properly")
	}

	if redisHost == "" || redisPort == "" {
		return Config{}, fmt.Errorf("redis environment variables are not set properly")
	}

	return Config{
		MySQLHost:     mysqlHost,
		MySQLPort:     mysqlPort,
		MySQLUser:     mysqlUser,
		MySQLPassword: mysqlPassword,
		MySQLDBName:   mysqlDBName,
		RedisHost:     redisHost,
		RedisPort:     redisPort,
	}, nil
}
