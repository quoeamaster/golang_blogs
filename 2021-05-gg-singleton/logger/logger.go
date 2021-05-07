package logger

import (
	"fmt"
	"sync"
	"time"
)

// Logger - class / object / structure to provide Log features.
// It is designed to be Stateful and hence could store Logger specific information.
type Logger struct{}

// Log - log the given [msg] to standard-out.
func (l *Logger) Log(msg string) {
	fmt.Printf("[%v] %v", time.Now().Format(time.RFC3339), msg)
}

// lock - the Mutex / Lock to guarantee only a code block is
// accessible ONLY by 1 thread / routine at a time.
var lock sync.Mutex

// singleton - the 1 and ONLY 1 instance of type Logger.
var singleton *Logger

// GetLogger - return the singleton instance.
func GetLogger() *Logger {
	lock.Lock()
	defer lock.Unlock()

	if singleton == nil {
		singleton = new(Logger)
	}
	return singleton
}
