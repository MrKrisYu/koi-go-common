package memory

import (
	"github.com/MrKrisYu/koi-go-common/config/loader"
	"github.com/MrKrisYu/koi-go-common/config/reader"
	"github.com/MrKrisYu/koi-go-common/config/source"
)

// WithSource appends a source to list of sources
func WithSource(s source.Source) loader.Option {
	return func(o *loader.Options) {
		o.Source = append(o.Source, s)
	}
}

// WithReader sets the config reader
func WithReader(r reader.Reader) loader.Option {
	return func(o *loader.Options) {
		o.Reader = r
	}
}
