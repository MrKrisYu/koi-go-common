package zap

import (
	"github.com/MrKrisYu/koi-go-common/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

type callerSkipKey struct{}

func WithCallerSkip(i int) logger.Option {
	return logger.SetOption(callerSkipKey{}, i)
}

type configKey struct{}

// WithConfig pass zap.Config to logger
func WithConfig(c zap.Config) logger.Option {
	return logger.SetOption(configKey{}, c)
}

type encoderConfigKey struct{}

// WithEncoderConfig pass zapcore.EncoderConfig to logger
func WithEncoderConfig(c zapcore.EncoderConfig) logger.Option {
	return logger.SetOption(encoderConfigKey{}, c)
}

type namespaceKey struct{}

func WithNamespace(namespace string) logger.Option {
	return logger.SetOption(namespaceKey{}, namespace)
}

type writerKey struct{}

type multiWriters struct {
	writers []io.Writer
}

func WithOutput(out ...io.Writer) logger.Option {
	return logger.SetOption(writerKey{}, multiWriters{writers: out})
}
