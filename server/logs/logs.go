package logs

import (
	"log"
	"os"
)

func InitFiles() {
	// Open log files
	logFile, err := os.OpenFile("./server/logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o664)
	if err != nil {
		log.Fatalf("Failed to open app log file: %v", err)
	}

	errorFile, err := os.OpenFile("./server/logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o664)
	if err != nil {
		log.Fatalf("Failed to open error log file: %v", err)
	}

	fatalFile, err := os.OpenFile("./server/logs/fatal.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o664)
	if err != nil {
		log.Fatalf("Failed to open fatal log file: %v", err)
	}

	// Setup loggers
	infoFILE = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorFILE = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	fatalFILE = log.New(fatalFile, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Set default log output for Fatal logs
	log.SetOutput(errorFile)
}

// GetLogger returns a logger with the specified prefix
func GetLogger(prefix string) *log.Logger {
	return errorFILE
}

var (
	infoFILE  *log.Logger
	errorFILE *log.Logger
	fatalFILE *log.Logger
)

func Fatalf(format string, v ...interface{}) {
	fatalFILE.Fatalf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	errorFILE.Printf(format, v...)
}

func Println(v ...interface{}) {
	infoFILE.Println(v...)
}

func Printf(format string, v ...interface{}) {
	infoFILE.Printf(format, v...)
}

func Print(v ...interface{}) {
	infoFILE.Print(v...)
}

func Fatal(v ...interface{}) {
	fatalFILE.Fatal(v...)
}

func Fatalln(v ...interface{}) {
	fatalFILE.Fatalln(v...)
}
