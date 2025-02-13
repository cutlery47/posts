package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type App struct {
	PostStorage
	UserStorage
}

type PostStorage struct {
	RestoreSource   string        `env:"RESTORE_SOURCE"`
	DumpDestination string        `env:"DUMP_DESTINATION"`
	DumpEnabled     bool          `env:"DUMP_ENABLED"`
	DumpInterval    time.Duration `env:"DUMP_INTERVAL"`
}

type UserStorage struct {
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
