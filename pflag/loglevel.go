package pflag

import (
	"fmt"

	logger "github.com/apsdehal/go-logger"
)

// LogLevel is a log level value.
type LogLevel logger.LogLevel

// LogLevelLitterals returns the litteral values of `LogLevel`.
var LogLevelLitterals = [...]string{
	"info",
	"debug",
	"notice",
	"warning",
	"error",
	"critical",
}

// Type implements `pflag.Value` interface.
func (s *LogLevel) Type() string {
	return "loglevel"
}

// Set implements `pflag.Value` interface.
func (s *LogLevel) Set(val string) error {
	switch val {
	case "info":
		*s = LogLevel(logger.InfoLevel)
	case "debug":
		*s = LogLevel(logger.DebugLevel)
	case "notice":
		*s = LogLevel(logger.NoticeLevel)
	case "warning":
		*s = LogLevel(logger.WarningLevel)
	case "error":
		*s = LogLevel(logger.ErrorLevel)
	case "critical":
		*s = LogLevel(logger.CriticalLevel)
	default:
		return fmt.Errorf("log level %s doesn't exists", s)
	}

	return nil
}

func (s *LogLevel) String() string {
	switch logger.LogLevel(*s) {
	case logger.InfoLevel:
		return "info"
	case logger.DebugLevel:
		return "debug"
	case logger.NoticeLevel:
		return "notice"
	case logger.WarningLevel:
		return "warning"
	case logger.ErrorLevel:
		return "error"
	case logger.CriticalLevel:
		return "critical"
	default:
		return ""
	}
}
