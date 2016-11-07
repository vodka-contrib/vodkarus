// Package vodkarus provides a middleware for vodka that logs request details
// via the logrus logging library
package vodkarus

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/insionng/vodka"
)

// New returns a new middleware handler with a default name and logger
func New() vodka.MiddlewareFunc {
	return NewWithName("vodkarus")
}

// NewWithName returns a new middleware handler with the specified name
func NewWithName(name string) vodka.MiddlewareFunc {
	return NewWithNameAndLogger(name, logrus.StandardLogger())
}

// NewWithNameAndLogger returns a new middleware handler with the specified name
// and logger
func NewWithNameAndLogger(name string, l *logrus.Logger) vodka.MiddlewareFunc {
	return func(next vodka.HandlerFunc) vodka.HandlerFunc {
		return func(c vodka.Context) error {
			start := time.Now()

			entry := l.WithFields(logrus.Fields{
				"request": c.Request().URI(),
				"method":  c.Request().Method(),
				"remote":  c.Request().RemoteAddress(),
			})

			if reqID := c.Request().Header().Get("X-Request-ID"); reqID != "" {
				entry = entry.WithField("request_id", reqID)
			}

			entry.Info("started handling request")

			if err := next(c); err != nil {
				c.Error(err)
			}

			latency := time.Since(start)

			entry = entry.WithFields(logrus.Fields{
				"status":      c.Response().Status(),
				"text_status": http.StatusText(c.Response().Status()),
				"took":        latency,
			})

			if c.Response().Status() == http.StatusNotFound {
				entry.Warn("completed handling request")
			} else {
				entry.Info("completed handling request")
			}

			return nil
		}
	}
}
