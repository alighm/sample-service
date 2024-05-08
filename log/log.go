package log

import (
	"context"
	stdlog "log"
	"os"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/alighm/sample-service/ctxutil"
)

var initialize sync.Once
var _logger logrus.FieldLogger

// RFC3339NanoFixed is time.RFC3339Nano with nanoseconds padded using zeros to
// ensure the formatted time is always the same number of characters.
const RFC3339NanoFixed = "2006-01-02T15:04:05.000000000Z07:00"

// GetLogger retrieves the current logger or creates a new one if needed.
func GetLogger(hooks ...logrus.Hook) logrus.FieldLogger {
	checkLogger(hooks...)
	return _logger
}

// SetLevel sets the level at which log messages are published/written.
func SetLevel(level string) {
	checkLogger()

	// If there's no explicit logging level specified, set the level to INFO
	if level == "" {
		level = "info"
	}

	loglevel, err := logrus.ParseLevel(level)
	if err == nil {
		// set default logger and the custom logger levels
		logrus.SetLevel(loglevel)
		_logger.(*logrus.Entry).Logger.SetLevel(loglevel)
	}
}

// SetService adds the Service field to each log message if name provided
func SetService(name string) {
	checkLogger()
	if name != "" {
		_logger = _logger.WithField("Service", name)
	}
}

// Logger configures a logger instance using the provided context
func Logger(ctx context.Context) logrus.FieldLogger {
	checkLogger()
	l := _logger

	if ctx != nil {
		if reqID := ctxutil.RequestID(ctx); reqID.String() != "" {
			l = l.WithField("CorrelationId", reqID.String())
		}
		if version := ctxutil.Version(ctx); version != "" {
			l = l.WithField("Version", version)
		}
	}
	return l
}

func checkLogger(hooks ...logrus.Hook) {
	if _logger == nil {
		newLogger(hooks...)
	}
}

func newLogger(hooks ...logrus.Hook) {
	l := logrus.New()
	// configure the default logger to include timestamps and quote empty fields to make visually
	// seeing an empty Field easier. These configurations will not impact or influence the
	// configuration of the logstash hook below.
	l.Formatter = &logrus.TextFormatter{
		TimestampFormat:  RFC3339NanoFixed,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
	}

	// Add each of the supplied hooks to the logger
	for _, h := range hooks {
		// add the hook to the logger
		l.Hooks.Add(h)
	}

	_logger = l.WithField("container", os.Getenv("HOSTNAME"))

	// configure standard logger to use the configured logger to write log entries
	initialize.Do(func() {
		// disable any flags that result in a prefix to the message
		// otherwise there will be duplicate timestamps, etc.
		stdlog.SetFlags(0)
		// use WriterLevel to define at what level library code messages should
		// be included into the logs. Most of the time the messages should be silent
		// unless additional diagnostics are required.
		stdlog.SetOutput(_logger.(*logrus.Entry).WriterLevel(logrus.DebugLevel))
	})

	_logger.Debug("Logger initialized")
}
