// Logging
package logger

import (
	"log"
	"time"
)

func Log(message string, level string) {
	log.SetFlags(0)
	log.Print(time.Now().Format("2006-01-02 15:04:05"), " [", level, "] ", message)
}
