package teavet

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Colors to represent different types of logs
var (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

// Log represents the structure of a log message
type Log struct {
	Filename  string    `json:"filename"`
	Line      string    `json:"line"`
	Timestamp time.Time `json:"timestamp"`

	Optional *Optional `json:"optional"`
}

// Optional represents the structure of the optional log fields
type Optional struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   error  `json:"error"`
}

// StartLog creates log file at the root of the current working directory
func StartLog() (err error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get working directory: %v", err)
		return err
	}

	// Create the log file
	file, err := os.OpenFile(path.Join(dir, "teavet.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Unable to create/append/open file: %v", err)
		return err
	}

	err = file.Close()
	if err != nil {
		log.Fatalf("Unable to close file: %v", err)
		return err
	}

	return nil
}

// FindLog finds a log file in the current working directory
func FindLog() (logFile string, err error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get working directory: %v", err)
		return "", err
	}

	// Look for the log file
	file, err := filepath.Glob(path.Join(dir, "teavet.log"))
	if err != nil {
		log.Fatalf("Unable to find the log file: %v", err)
		return "", err
	}

	logFile = strings.Join(file, "")

	return logFile, nil
}

// BuildLog formats and writes content to the log file
func BuildLog(optional *Optional) (logContent *Log, err error) {
	// Get information about the file and the line
	_, filename, line, _ := runtime.Caller(1)

	logContent = &Log{
		Filename:  filename,
		Line:      strconv.Itoa(line),
		Timestamp: time.Now().UTC(),
		Optional:  optional,
	}

	// Find the log file
	logFile, err := FindLog()
	if err != nil {
		log.Fatalf("Unable to find the log file: %v", err)
		return logContent, err
	}

	// Format and save log file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Unable to create/append/open file: %v", err)
		return logContent, err
	}

	formatLog, err := json.MarshalIndent(logContent, "", " ")
	if err != nil {
		log.Fatalf("Unable to marshal log: %v", err)
		return logContent, err
	}

	_, err = file.Write([]byte(formatLog))
	if err != nil {
		log.Fatalf("Unable to write to file: %v", err)
		return logContent, err
	}

	err = file.Close()
	if err != nil {
		log.Fatalf("Unable to close file: %v", err)
		return logContent, err
	}

	return logContent, nil
}

// Logger prints a log with filename, line, timestamp, status, message and error
func Logger(logStatus, logMessage string, logError error) {
	logContent, err := BuildLog(&Optional{Status: logStatus, Message: logMessage, Error: logError})
	if err != nil {
		log.Fatalf("Unable to build log: %v", err)
	}

	formatLog := fmt.Sprintf("\n### Filename: %s\n### Line: %s\n### Timestamp: %s\n### Status: %s\n### Message: %s\n### Error: %v\n", logContent.Filename, logContent.Line, logContent.Timestamp, logContent.Optional.Status, logContent.Optional.Message, logContent.Optional.Error)

	// Remove the timestamp provided by log package
	log.SetFlags(0)

	switch strings.ToLower(logContent.Optional.Status) {
	case "success":
		log.Printf(Green+"\n%v\n"+Reset, formatLog)
	case "info":
		log.Printf(Cyan+"\n%v\n"+Reset, formatLog)
	case "warning":
		log.Printf(Yellow+"\n%v\n"+Reset, formatLog)
	case "fatal":
		log.Fatalf(Purple+"\n%v\n"+Reset, formatLog)
	case "panic":
		log.Panicf(Red+"\n%v\n"+Reset, formatLog)
	default:
		log.Printf(Blue+"\n%v\n"+Reset, formatLog)
	}
}

// SimpleLogger prints a log with filename, line and timestamp
func SimpleLogger() {
	logContent, err := BuildLog(&Optional{})
	if err != nil {
		log.Fatalf("Unable to build log: %v", err)
	}

	formatLog := fmt.Sprintf("\n### Filename: %s\n### Line: %s\n### Timestamp: %s\n", logContent.Filename, logContent.Line, logContent.Timestamp)

	// Remove the timestamp provided by log package
	log.SetFlags(0)

	log.Printf(Blue+"\n%v\n"+Reset, formatLog)
}

// StatusLogger prints a log with filename, line, timestamp and status
func StatusLogger(logStatus string) {
	logContent, err := BuildLog(&Optional{Status: logStatus})
	if err != nil {
		log.Fatalf("Unable to build log: %v", err)
	}

	formatLog := fmt.Sprintf("\n### Filename: %s\n### Line: %s\n### Timestamp: %s\n### Status: %s\n", logContent.Filename, logContent.Line, logContent.Timestamp, logContent.Optional.Status)

	// Remove the timestamp provided by log package
	log.SetFlags(0)

	switch strings.ToLower(logContent.Optional.Status) {
	case "success":
		log.Printf(Green+"\n%v\n"+Reset, formatLog)
	case "info":
		log.Printf(Cyan+"\n%v\n"+Reset, formatLog)
	case "warning":
		log.Printf(Yellow+"\n%v\n"+Reset, formatLog)
	case "fatal":
		log.Fatalf(Purple+"\n%v\n"+Reset, formatLog)
	case "panic":
		log.Panicf(Red+"\n%v\n"+Reset, formatLog)
	default:
		log.Printf(Blue+"\n%v\n"+Reset, formatLog)
	}
}

// MessageLogger prints a log with filename, line, timestamp and message
func MessageLogger(logMessage string) {
	logContent, err := BuildLog(&Optional{Message: logMessage})
	if err != nil {
		log.Fatalf("Unable to build log: %v", err)
	}

	formatLog := fmt.Sprintf("\n### Filename: %s\n### Line: %s\n### Timestamp: %s\n### Message: %s\n", logContent.Filename, logContent.Line, logContent.Timestamp, logContent.Optional.Message)

	// Remove the timestamp provided by log package
	log.SetFlags(0)

	log.Printf(Blue+"\n%v\n"+Reset, formatLog)
}

// ErrorLogger prints a log with filename, line, timestamp and error
func ErrorLogger(logError error) {
	logContent, err := BuildLog(&Optional{Error: logError})
	if err != nil {
		log.Fatalf("Unable to build log: %v", err)
	}

	formatLog := fmt.Sprintf("\n### Filename: %s\n### Line: %s\n### Timestamp: %s\n### Error: %v\n", logContent.Filename, logContent.Line, logContent.Timestamp, logContent.Optional.Error)

	// Remove the timestamp provided by log package
	log.SetFlags(0)

	log.Fatalf(Purple+"\n%v\n"+Reset, formatLog)
}

// PanicLogger prints a log with filename, line, timestamp and error
func PanicLogger(logError error) {
	logContent, err := BuildLog(&Optional{Error: logError})
	if err != nil {
		log.Fatalf("Unable to build log: %v", err)
	}

	formatLog := fmt.Sprintf("\n### Filename: %s\n### Line: %s\n### Timestamp: %s\n### Error: %v\n", logContent.Filename, logContent.Line, logContent.Timestamp, logContent.Optional.Error)

	// Remove the timestamp provided by log package
	log.SetFlags(0)

	log.Panicf(Red+"\n%v\n"+Reset, formatLog)
}
