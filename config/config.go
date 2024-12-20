package config

import "os"

type Config struct {
	DbUrl string
}

func LoadConfig() Config {
	return Config{
		DbUrl: os.Getenv("DATABASE_URL"),
	}
}
