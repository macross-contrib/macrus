// Package macrus provides a middleware for macross that logs request details
// via the logrus logging library
package macrus

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/insionng/macross"
)

// New returns a new middleware handler with a default name and logger
func New() macross.Handler {
	return NewWithName("macrus")
}

// NewWithName returns a new middleware handler with the specified name
func NewWithName(name string) macross.Handler {
	return NewWithNameAndLogger(name, logrus.StandardLogger())
}

// NewWithNameAndLogger returns a new middleware handler with the specified name
// and logger
func NewWithNameAndLogger(name string, l *logrus.Logger) macross.Handler {
	return func(c *macross.Context) error {
		start := time.Now()

		entry := l.WithFields(logrus.Fields{
			"request": string(c.Request.URI().Path()),
			"method":  string(c.Request.Header.Method()),
			"remote":  c.RemoteAddr().String(),
		})

		if reqID := string(c.Request.Header.Peek("X-Request-ID")); reqID != "" {
			entry = entry.WithField("request_id", reqID)
		}

		entry.Info("started handling request")

		if err := c.Next(); err != nil {
			c.Error(err.Error(), macross.StatusInternalServerError)
		}

		latency := time.Since(start)

		entry = entry.WithFields(logrus.Fields{
			"status":      c.Response.StatusCode(),
			"text_status": macross.StatusText(c.Response.StatusCode()),
			"took":        latency,
		})

		if c.Response.StatusCode() == http.StatusNotFound {
			entry.Warn("completed handling request")
		} else {
			entry.Info("completed handling request")
		}

		return nil
	}
}
