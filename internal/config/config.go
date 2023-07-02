package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"golang.org/x/exp/slog"
)

type Config struct {
	Logger      Logger
	Server      Server
	GRPC        GRPC
	Postgres    Postgres
	FileService FileService
}

type Logger struct {
	Level slog.Level `env:"LOGGER_LEVEL" env-default:"debug"`
}

type Server struct {
	Addr string `env:"SERVER_ADDR"`
}

type GRPC struct {
	Addr string `env:"GRPC_ADDR"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DB       string `env:"POSTGRES_DB"`
}

type FileService struct {
	MaxFileSize int64 `env:"FILE_SERVICE_MAX_FILE_SIZE"`
}

func New() (Config, error) {
	var config Config
	err := cleanenv.ReadEnv(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
