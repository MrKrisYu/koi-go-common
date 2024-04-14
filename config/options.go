package config

import (
	"github.com/MrKrisYu/koi-go-common/config/encoder"
	"github.com/MrKrisYu/koi-go-common/config/encoder/yaml"
	"github.com/MrKrisYu/koi-go-common/config/entity"
)

const (
	defaultConfigPath = "./application.yaml"
)

var (
	defaultEncoder = yaml.NewEncoder()
)

type Options struct {
	// source path
	ConfigPath string
	// Encoder
	Encoder encoder.Encoder
	// for alternative data
	Entities []entity.Entity
	// for callbacks
	Callbacks []func()
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	// default options
	options := Options{
		ConfigPath: defaultConfigPath,
		Encoder:    defaultEncoder,
		Entities:   []entity.Entity{},
		Callbacks:  []func(){},
	}
	// configures options
	for _, o := range opts {
		o(&options)
	}

	return options
}

func WithPath(path string) Option {
	return func(o *Options) {
		o.ConfigPath = path
	}
}

func WithEncoder(e encoder.Encoder) Option {
	return func(o *Options) {
		o.Encoder = e
	}
}

func WithEntity(e ...entity.Entity) Option {
	return func(o *Options) {
		o.Entities = append(o.Entities, e...)
	}
}

func WithCallback(f ...func()) Option {
	return func(o *Options) {
		o.Callbacks = append(o.Callbacks, f...)
	}
}
