// Package logger implements github.com/the-anna-project/logger.Service. This
// logger interface is to simply log output to gather runtime information.
package logger

import (
	"os"
	"time"

	kitlog "github.com/go-kit/kit/log"

	microerror "github.com/giantswarm/microkit/error"
)

// Config represents the configuration used to create a new logger service.
type Config struct {
	// Settings.
	TimestampFormatter kitlog.Valuer
}

// DefaultConfig provides a default configuration to create a new logger service
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Settings.
		TimestampFormatter: func() interface{} {
			return time.Now().UTC().Format("06-01-02 15:04:05.000")
		},
	}
}

// New creates a new configured logger service.
func New(config Config) (Service, error) {
	// Settings.
	if config.TimestampFormatter == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "timestamp formatter must not be empty")
	}

	kitLogger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	kitLogger = kitlog.NewContext(kitLogger).With(
		"ts", config.TimestampFormatter,
		"caller", kitlog.DefaultCaller,
	)

	newLogger := &service{
		logger: kitLogger,
	}

	return newLogger, nil
}

type service struct {
	logger Service
}

func (s *service) Log(v ...interface{}) error {
	return s.logger.Log(v...)
}
