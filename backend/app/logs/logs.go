package logs

import (
	"log"
	"os"
)

var (
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	FatalLog *log.Logger
	Err      error
)

func init() {
	// Ensure the log directory exists
	logDir := "./server/logs"
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Open log files
	logFile, err := os.OpenFile(logDir+"/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o664)
	if err != nil {
		log.Fatalf("Failed to open app log file: %v", err)
	}

	errorFile, err := os.OpenFile(logDir+"/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o664)
	if err != nil {
		log.Fatalf("Failed to open error log file: %v", err)
	}

	fatalFile, err := os.OpenFile(logDir+"/fatal.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o664)
	if err != nil {
		log.Fatalf("Failed to open fatal log file: %v", err)
	}

	// Setup loggers
	InfoLog = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	FatalLog = log.New(fatalFile, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Set default log output for Fatal logs
	log.SetOutput(errorFile)
}
