package loader

import (
	"context"
	"github.com/MrKrisYu/koi-go-common/config/reader"
	"github.com/MrKrisYu/koi-go-common/config/source"
)

type Option func(o *Options)

type Options struct {
	Reader reader.Reader
	Source []source.Source

	// for alternative data
	Context context.Context
}
