package config

import (
	"fmt"
	"github.com/MrKrisYu/koi-go-common/config/entity"
	"os"
)

var (
	Global *Config
)

type Config struct {
	// encoder
	Opts Options
	// config data
	Settings *entity.Settings
	// for alternative data
	entities map[string]entity.Entity
}

func Setup(option ...Option) {
	options := NewOptions(option...)
	Global = &Config{
		Opts:     options,
		Settings: entity.NewSettings(),
		entities: make(map[string]entity.Entity),
	}
	bytes, err := os.ReadFile(Global.Opts.ConfigPath)
	if err != nil {
		panic(err)
	}
	err = Global.Opts.Encoder.Decode(bytes, Global.Settings)
	if err != nil {
		panic(err)
	}
	// initialize settings
	Global.Settings.Init()
	WithEntity(Global.Settings.Root.ApplicationConfig, Global.Settings.Root.LogConfig)(&Global.Opts)

	// handle extra settings
	for _, entry := range Global.Opts.Entities {
		Global.entities[entry.Key()] = entry
	}
	for _, callback := range Global.Opts.Callbacks {
		callback()
	}
}

func GetEntity[T any](config *Config, key string) (T, error) {
	e, exist := config.entities[key]
	if !exist {
		return *new(T), fmt.Errorf("cannot find entity")
	}
	result, ok := e.(T)
	if !ok {
		t := new(T)
		return *t, fmt.Errorf("cannot convert %T to %T", e, *t)
	}
	return result, nil
}
