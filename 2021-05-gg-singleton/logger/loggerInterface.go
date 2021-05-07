package logger

import (
	"fmt"
	"time"
)

// ILog - logger interface.
type ILog interface {
	// Log - logging the given [msg] to standard out.
	Log(msg string)
}

// loggerImpl - a private implemenation of [ILog] interface.
type loggerImpl struct{}

func (l *loggerImpl) Log(msg string) {
	fmt.Printf("[%v] %v", time.Now().Format(time.RFC3339), msg)
}

var singletonImpl ILog

// GetLoggerImpl - return the [ILog] implementation
// instead of the underlying [loggerImpl].
// Now nobody can create an instance of loggerImpl by the 'new' operator.
func GetLoggerImpl() ILog {
	lock.Lock()
	defer lock.Unlock()

	if singletonImpl == nil {
		singletonImpl = new(loggerImpl)
	}
	return singletonImpl
}
