package logs

import (
	"log"
	"os"
)

var logger *log.Logger

// InitLogger initializes a logger
func InitLogger() {
	logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// LogInfo logs an info message
func LogInfo(message string) {
	logger.Println(message)
}

// LogError logs an error message
func LogError(err error) {
	logger.Printf("ERROR: %v", err)
}
