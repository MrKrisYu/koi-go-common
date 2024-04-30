package config

import (
	"github.com/MrKrisYu/koi-go-common/config/reader"
	"github.com/MrKrisYu/koi-go-common/config/source"
)

// Config is an interface abstraction for dynamic configuration
type Config interface {
	// Values provide the reader.Values interface
	reader.Values
	// Init the DefaultConfig
	Init(opts ...Option) error
	// Options in the DefaultConfig
	Options() Options
	// Close Stop the DefaultConfig loader/watcher
	Close() error
	// Load DefaultConfig sources
	Load(source ...source.Source) error
	// Sync Force a source changeset sync
	Sync() error
}
