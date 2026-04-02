package log

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	instance *log.Logger
	once     sync.Once
)

// Get returns the global logger instance
func Get() *log.Logger {
	once.Do(func() {
		instance = log.New(os.Stdout, "[FABRIK] ", log.LstdFlags|log.Lshortfile)
	})
	return instance
}

// Helper functions for convenience
func Info(msg string, args ...interface{}) {
	Get().Output(2, fmt.Sprintf("[INFO] "+msg, args...))
}

func Error(msg string, args ...interface{}) {
	Get().Output(2, fmt.Sprintf("[ERROR] "+msg, args...))
}

func Debug(msg string, args ...interface{}) {
	Get().Output(2, fmt.Sprintf("[DEBUG] "+msg, args...))
}
