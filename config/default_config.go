package config

import (
	"github.com/MrKrisYu/koi-go-common/config/loader"
	"github.com/MrKrisYu/koi-go-common/config/loader/memory"
	"github.com/MrKrisYu/koi-go-common/config/reader"
	"github.com/MrKrisYu/koi-go-common/config/reader/json"
	"github.com/MrKrisYu/koi-go-common/config/source"
	"sync"
)

type DefaultConfig struct {
	exit chan bool
	opts Options

	sync.RWMutex
	// the current snapshot
	snap *loader.Snapshot
	// the current values
	vals reader.Values
}

func NewConfig(opts ...Option) (Config, error) {
	var c DefaultConfig

	err := c.Init(opts...)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *DefaultConfig) Init(opts ...Option) error {
	c.opts = Options{
		Reader: json.NewReader(),
	}
	c.exit = make(chan bool)
	for _, o := range opts {
		o(&c.opts)
	}

	// default loader uses the configured reader
	if c.opts.Loader == nil {
		c.opts.Loader = memory.NewLoader(memory.WithReader(c.opts.Reader))
	}

	err := c.opts.Loader.Load(c.opts.Source...)
	if err != nil {
		return err
	}

	c.snap, err = c.opts.Loader.Snapshot()
	if err != nil {
		return err
	}

	c.vals, err = c.opts.Reader.Values(c.snap.ChangeSet)
	if err != nil {
		return err
	}
	if c.opts.Entity != nil {
		_ = c.vals.Scan(c.opts.Entity)
	}

	return nil
}

func (c *DefaultConfig) Options() Options {
	return c.opts
}

func (c *DefaultConfig) Map() map[string]interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.vals.Map()
}

func (c *DefaultConfig) Scan(v interface{}) error {
	c.RLock()
	defer c.RUnlock()
	return c.vals.Scan(v)
}

// Sync loads all the sources, calls the parser and updates the DefaultConfig
func (c *DefaultConfig) Sync() error {
	if err := c.opts.Loader.Sync(); err != nil {
		return err
	}

	snap, err := c.opts.Loader.Snapshot()
	if err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()

	c.snap = snap
	vals, err := c.opts.Reader.Values(snap.ChangeSet)
	if err != nil {
		return err
	}
	c.vals = vals

	return nil
}

func (c *DefaultConfig) Close() error {
	select {
	case <-c.exit:
		return nil
	default:
		close(c.exit)
	}
	return nil
}

func (c *DefaultConfig) Get(path ...string) reader.Value {
	c.RLock()
	defer c.RUnlock()

	// did sync actually work?
	if c.vals != nil {
		return c.vals.Get(path...)
	}

	// no value
	return newValue()
}

func (c *DefaultConfig) Set(val interface{}, path ...string) {
	c.Lock()
	defer c.Unlock()

	if c.vals != nil {
		c.vals.Set(val, path...)
	}

	return
}

func (c *DefaultConfig) Del(path ...string) {
	c.Lock()
	defer c.Unlock()

	if c.vals != nil {
		c.vals.Del(path...)
	}

	return
}

func (c *DefaultConfig) Bytes() []byte {
	c.RLock()
	defer c.RUnlock()

	if c.vals == nil {
		return []byte{}
	}

	return c.vals.Bytes()
}

func (c *DefaultConfig) Load(sources ...source.Source) error {
	if err := c.opts.Loader.Load(sources...); err != nil {
		return err
	}

	snap, err := c.opts.Loader.Snapshot()
	if err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()

	c.snap = snap
	vals, err := c.opts.Reader.Values(snap.ChangeSet)
	if err != nil {
		return err
	}
	c.vals = vals

	return nil
}

func (c *DefaultConfig) String() string {
	return "DefaultConfig"
}
