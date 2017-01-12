// Package logger implements github.com/the-anna-project/logger.Service. This
// logger interface is to simply log output to gather runtime information.
package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-stack/stack"
)

// ServiceConfig represents the configuration used to create a new service.
type ServiceConfig struct {
	// Settings.
	Caller             kitlog.Valuer
	IOWriter           io.Writer
	TimestampFormatter kitlog.Valuer
}

// DefaultServiceConfig provides a default configuration to create a new service
// by best effort.
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		// Settings.
		Caller: func() interface{} {
			return fmt.Sprintf("%+v", stack.Caller(4))
		},
		IOWriter: ioutil.Discard,
		TimestampFormatter: func() interface{} {
			return time.Now().UTC().Format("06-01-02 15:04:05.000")
		},
	}
}

// NewService creates a new configured service.
func NewService(config ServiceConfig) (Service, error) {
	// Settings.
	if config.Caller == nil {
		return nil, maskAnyf(invalidConfigError, "caller must not be empty")
	}
	if config.TimestampFormatter == nil {
		return nil, maskAnyf(invalidConfigError, "timestamp formatter must not be empty")
	}

	kitLogger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(config.IOWriter))
	kitLogger = kitlog.NewContext(kitLogger).With(
		"caller", config.Caller,
		"time", config.TimestampFormatter,
	)

	newService := &service{
		logger: kitLogger,
	}

	return newService, nil
}

type service struct {
	logger Service
}

func (s *service) Log(v ...interface{}) error {
	return s.logger.Log(v...)
}
