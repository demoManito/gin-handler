package handler

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

// Level type
type Level int

// logging levels
const (
	PanicLevel = Level(logrus.PanicLevel)
	FatalLevel = Level(logrus.FatalLevel)
	ErrorLevel = Level(logrus.ErrorLevel)
	WarnLevel  = Level(logrus.WarnLevel)
	InfoLevel  = Level(logrus.InfoLevel)
	DebugLevel = Level(logrus.DebugLevel)
	TraceLevel = Level(logrus.TraceLevel)
)

// ContextLogger ContextLogger interface
type ContextLogger interface {
	Log(Level, string)
}

// ContextError implements the error interface.
type ContextError struct {
	title  string
	caller *runtime.Frame
	err    error
	data   logrus.Fields
}

// NewContextError new ContextError with title and default code
func NewContextError(title string) *ContextError {
	return &ContextError{title: title}
}

// Error returns the error info
func (e *ContextError) Error() string {
	if e.err == nil {
		return e.Title()
	}
	return e.Title() + ": " + e.err.Error()
}

// Title gets the title of the ContextError
func (e *ContextError) Title() string {
	return e.title
}
