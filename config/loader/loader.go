package loader

import (
	"github.com/MrKrisYu/koi-go-common/config/source"
)

// Loader manages loading sources
type Loader interface {
	// Close Stop the loader
	Close() error
	// Load the sources
	Load(...source.Source) error
	// Snapshot A Snapshot of loaded config
	Snapshot() (*Snapshot, error)
	// Sync Force sync of sources
	Sync() error
	// String Name of loader
	String() string
}

// Snapshot is a merged ChangeSet
type Snapshot struct {
	// The merged ChangeSet
	ChangeSet *source.ChangeSet
	// Version Deterministic and comparable version of the snapshot
	Version string
}

// Copy snapshot
func Copy(s *Snapshot) *Snapshot {
	cs := *(s.ChangeSet)

	return &Snapshot{
		ChangeSet: &cs,
		Version:   s.Version,
	}
}
