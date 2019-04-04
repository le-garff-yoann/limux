package processor

import (
	logger "github.com/apsdehal/go-logger"
)

// Processor is the common interface of `In`, `Out`, ....
type Processor interface {
	Configure() error
	Start(chan (Event)) error
	Stop() error
	String() string
}

const (
	// Log is a type of `Event`.
	Log = iota
	// Bgn is a type of `Event`.
	Bgn
	// Fin is a type of `Event`.
	Fin
)

// Event basically represent (atm) a log message.
type Event struct {
	Processor Processor       `json:"processor"`
	Type      int             `json:"type"`
	Message   string          `json:"message"`
	Level     logger.LogLevel `json:"level"`
}

// EnvironmentError could typically be used into a validate subcommand....
type EnvironmentError string

func (s EnvironmentError) Error() string {
	return string(s)
}
