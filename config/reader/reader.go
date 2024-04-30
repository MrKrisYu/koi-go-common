package reader

import "github.com/MrKrisYu/koi-go-common/config/source"

// Reader is an interface for merging changesets
type Reader interface {
	Merge(...*source.ChangeSet) (*source.ChangeSet, error)
	Values(*source.ChangeSet) (Values, error)
	String() string
}
