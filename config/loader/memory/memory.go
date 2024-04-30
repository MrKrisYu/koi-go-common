package memory

import (
	"errors"
	"fmt"
	"github.com/MrKrisYu/koi-go-common/config/loader"
	"github.com/MrKrisYu/koi-go-common/config/reader"
	"github.com/MrKrisYu/koi-go-common/config/reader/json"
	"github.com/MrKrisYu/koi-go-common/config/source"
	"strings"
	"sync"
	"time"
)

type memory struct {
	exit chan bool
	opts loader.Options

	sync.RWMutex
	// the current snapshot
	snap *loader.Snapshot
	// the current values
	vals reader.Values
	// all the sets
	sets []*source.ChangeSet
	// all the sources
	sources []source.Source
}

type updateValue struct {
	version string
	value   reader.Value
}

func (m *memory) loaded() bool {
	var loaded bool
	m.RLock()
	if m.vals != nil {
		loaded = true
	}
	m.RUnlock()
	return loaded
}

// reload reads the sets and creates new values
func (m *memory) reload() error {
	m.Lock()

	// merge sets
	set, err := m.opts.Reader.Merge(m.sets...)
	if err != nil {
		m.Unlock()
		return err
	}

	// set values
	m.vals, _ = m.opts.Reader.Values(set)
	m.snap = &loader.Snapshot{
		ChangeSet: set,
		Version:   getVersion(),
	}

	m.Unlock()

	return nil
}

// Snapshot returns a snapshot of the current loaded config
func (m *memory) Snapshot() (*loader.Snapshot, error) {
	if m.loaded() {
		m.RLock()
		snap := loader.Copy(m.snap)
		m.RUnlock()
		return snap, nil
	}

	// not loaded, sync
	if err := m.Sync(); err != nil {
		return nil, err
	}

	// make copy
	m.RLock()
	snap := loader.Copy(m.snap)
	m.RUnlock()

	return snap, nil
}

// Sync loads all the sources, calls the parser and updates the config
func (m *memory) Sync() error {
	//nolint:prealloc
	var sets []*source.ChangeSet

	m.Lock()

	// read the oneSource
	var gerr []string

	for _, oneSource := range m.sources {
		ch, err := oneSource.Read()
		if err != nil {
			gerr = append(gerr, err.Error())
			continue
		}
		sets = append(sets, ch)
	}

	// merge sets
	set, err := m.opts.Reader.Merge(sets...)
	if err != nil {
		m.Unlock()
		return err
	}

	// set values
	vals, err := m.opts.Reader.Values(set)
	if err != nil {
		m.Unlock()
		return err
	}
	m.vals = vals
	m.snap = &loader.Snapshot{
		ChangeSet: set,
		Version:   getVersion(),
	}

	m.Unlock()

	if len(gerr) > 0 {
		return fmt.Errorf("oneSource loading errors: %s", strings.Join(gerr, "\n"))
	}

	return nil
}

func (m *memory) Close() error {
	select {
	case <-m.exit:
		return nil
	default:
		close(m.exit)
	}
	return nil
}

func (m *memory) Get(path ...string) (reader.Value, error) {
	if !m.loaded() {
		if err := m.Sync(); err != nil {
			return nil, err
		}
	}

	m.Lock()
	defer m.Unlock()

	// did sync actually work?
	if m.vals != nil {
		return m.vals.Get(path...), nil
	}

	// assuming vals is nil
	// create new vals

	ch := m.snap.ChangeSet

	// we are truly screwed, trying to load in a hacked way
	v, err := m.opts.Reader.Values(ch)
	if err != nil {
		return nil, err
	}

	// lets set it just because
	m.vals = v

	if m.vals != nil {
		return m.vals.Get(path...), nil
	}

	// ok we're going hardcore now

	return nil, errors.New("no values")
}

func (m *memory) Load(sources ...source.Source) error {
	var gerrors []string

	for _, oneSource := range sources {
		set, err := oneSource.Read()
		if err != nil {
			gerrors = append(gerrors,
				fmt.Sprintf("error loading oneSource %s: %v",
					oneSource,
					err))
			// continue processing
			continue
		}
		m.Lock()
		m.sources = append(m.sources, oneSource)
		m.sets = append(m.sets, set)
		m.Unlock()
	}

	if err := m.reload(); err != nil {
		gerrors = append(gerrors, err.Error())
	}

	// Return errors
	if len(gerrors) != 0 {
		return errors.New(strings.Join(gerrors, "\n"))
	}
	return nil
}

func (m *memory) String() string {
	return "memory"
}

func getVersion() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func NewLoader(opts ...loader.Option) loader.Loader {
	options := loader.Options{
		Reader: json.NewReader(),
	}

	for _, o := range opts {
		o(&options)
	}

	m := &memory{
		exit:    make(chan bool),
		opts:    options,
		sources: options.Source,
	}

	m.sets = make([]*source.ChangeSet, len(options.Source))

	for i, s := range options.Source {
		m.sets[i] = &source.ChangeSet{Source: s.String()}
	}

	return m
}
