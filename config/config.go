package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type App struct {
	Handler
	Service
	Storage
	HTTPServer
}

type Handler struct {
}

type Service struct {
}

type Storage struct {
	PostStorage
	UserStorage
}

type PostStorage struct {
	RestoreSource   string        `env:"RESTORE_SOURCE" env-default:"dump"`
	DumpDestination string        `env:"DUMP_DESTINATION" env-default:"dump"`
	DumpEnabled     bool          `env:"DUMP_ENABLED" env-default:"true"`
	DumpInterval    time.Duration `env:"DUMP_INTERVAL" env-default:"5s"`
}

type UserStorage struct {
	SessionDuration time.Duration `env:"SESSION_DURATION" env-default:"24h"`
}

type HTTPServer struct {
	BindAddress     string        `env:"BIND_ADDRESS" env-default:"localhost"`
	BindPort        string        `env:"BIND_PORT" env-default:"8000"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"5s"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT" env-default:"5s"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" env-default:"5s"`
}

func New(env string) (*App, error) {
	conf := &App{}

	if err := godotenv.Overload(env); err != nil {
		return nil, fmt.Errorf("godotenv.Overload: %v", err)
	}

	if err := cleanenv.ReadEnv(conf); err != nil {
		return nil, fmt.Errorf("cleanenv.Readenv: %v", err)
	}

	return conf, nil
}
