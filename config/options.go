package config

import (
	"context"
	"github.com/MrKrisYu/koi-go-common/config/loader"
	"github.com/MrKrisYu/koi-go-common/config/reader"
	"github.com/MrKrisYu/koi-go-common/config/source"
)

type Option func(o *Options)

type Options struct {
	Loader loader.Loader
	Reader reader.Reader
	Source []source.Source

	// for alternative data
	Context context.Context

	Entity Entity
}

// WithLoader sets the loader for manager DefaultConfig
func WithLoader(l loader.Loader) Option {
	return func(o *Options) {
		o.Loader = l
	}
}

// WithSource appends a source to list of sources
func WithSource(s source.Source) Option {
	return func(o *Options) {
		o.Source = append(o.Source, s)
	}
}

// WithReader sets the DefaultConfig reader
func WithReader(r reader.Reader) Option {
	return func(o *Options) {
		o.Reader = r
	}
}

// WithEntity sets the DefaultConfig Entity
func WithEntity(e Entity) Option {
	return func(o *Options) {
		o.Entity = e
	}
}
