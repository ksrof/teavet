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

// 8-16 ANSI color codes for normal and bright
var (
	Reset         = "\u001b[0m"
	BrightBlack   = "\u001b[30;1m"
	BrightRed     = "\u001b[31;1m"
	BrightGreen   = "\u001b[32;1m"
	BrightYellow  = "\u001b[33;1m"
	BrightBlue    = "\u001b[34;1m"
	BrightMagenta = "\u001b[35;1m"
	BrightCyan    = "\u001b[36;1m"
	BrightWhite   = "\u001b[37;1m"
)

// LoggerContent represents the structure of the log.
type LoggerContent struct {
	Filename  string `json:"filename"`
	Line      string `json:"line"`
	Timestamp string `json:"timestamp"`

	*LoggerStatus  `json:"logger_status,omitempty"`
	*LoggerMessage `json:"logger_message,omitempty"`
	*LoggerFault   `json:"logger_fault,omitempty"`
}

// LoggerStatus represents the log status field.
type LoggerStatus struct {
	Status string `json:"status,omitempty"`
}

// LoggerMessage represents the log message field.
type LoggerMessage struct {
	Message string `json:"message,omitempty"`
}

// LoggerFault represents the log fault field.
type LoggerFault struct {
	Fault string `json:"fault,omitempty"`
}

// StartLogger creates a log file at the root of the current directory.
func StartLogger() error {
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

// FindLogger looks for a log file in the current directory.
func FindLogger() (string, error) {
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

	loggerFile := strings.Join(file, "")

	return loggerFile, nil
}

// SaveLogger formats and writes content to the log file.
func SaveLogger(logger *LoggerContent) error {
	// Find the log file
	loggerFile, err := FindLogger()
	if err != nil {
		log.Fatalf("Unable to find the log file: %v", err)
		return err
	}

	// Format and save log file
	file, err := os.OpenFile(loggerFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Unable to create/append/open file: %v", err)
		return err
	}

	formated, err := json.MarshalIndent(logger, "", " ")
	if err != nil {
		log.Fatalf("Unable to marshal log: %v", err)
		return err
	}

	_, err = file.Write([]byte(formated))
	if err != nil {
		log.Fatalf("Unable to write to file: %v", err)
		return err
	}

	err = file.Close()
	if err != nil {
		log.Fatalf("Unable to close file: %v", err)
		return err
	}

	return nil
}

// FormatLogger formats the log output.
func FormatLogger(logger *LoggerContent, color string) string {
	if len(logger.LoggerStatus.Status) > 0 {
		formated := fmt.Sprintf("%s\n| Filename: %s\n| Line: %s\n| Timestamp: %s\n| Status: %s\n%s", color, logger.Filename, logger.Line, logger.Timestamp, logger.LoggerStatus.Status, Reset)
		return formated
	}

	if len(logger.LoggerMessage.Message) > 0 {
		formated := fmt.Sprintf("%s\n| Filename: %s\n| Line: %s\n| Timestamp: %s\n| Message: %s\n%s", color, logger.Filename, logger.Line, logger.Timestamp, logger.LoggerMessage.Message, Reset)
		return formated
	}

	if len(logger.LoggerFault.Fault) > 0 {
		formated := fmt.Sprintf("%s\n| Filename: %s\n| Line: %s\n| Timestamp: %s\n| Fault: %v\n%s", color, logger.Filename, logger.Line, logger.Timestamp, logger.LoggerFault.Fault, Reset)
		return formated
	}

	if len(logger.LoggerStatus.Status) <= 0 && len(logger.LoggerMessage.Message) <= 0 && len(logger.LoggerFault.Fault) <= 0 {
		formated := fmt.Sprintf("%s\n| Filename: %s\n| Line: %s\n| Timestamp: %s\n%s", color, logger.Filename, logger.Line, logger.Timestamp, Reset)
		return formated
	}

	formated := fmt.Sprintf("%s\n| Filename: %s\n| Line: %s\n| Timestamp: %s\n| Status: %s\n| Message: %s\n| Fault: %v\n%s", color, logger.Filename, logger.Line, logger.Timestamp, logger.LoggerStatus.Status, logger.LoggerMessage.Message, logger.LoggerFault.Fault, Reset)
	return formated
}

// Complete outputs a log with a given status, message and fault.
func Complete(status, message string, fault error) {
	// Disable default log timestamp
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Get information about the file and the line
	_, filename, line, _ := runtime.Caller(1)

	logger := &LoggerContent{
		Filename:      filename,
		Line:          strconv.Itoa(line),
		Timestamp:     time.Now().Format(time.RFC3339),
		LoggerStatus:  &LoggerStatus{Status: status},
		LoggerMessage: &LoggerMessage{Message: message},
		LoggerFault:   &LoggerFault{Fault: fmt.Sprint(fault)},
	}

	err := SaveLogger(logger)
	if err != nil {
		log.Fatalf("Unable to save logger to log file: %v", err)
	}

	switch logger.LoggerStatus.Status {
	case "success":
		log.Print(FormatLogger(logger, BrightGreen))
	case "info":
		log.Print(FormatLogger(logger, BrightBlue))
	case "error":
		log.Print(FormatLogger(logger, BrightYellow))
	case "fatal":
		log.Fatal(FormatLogger(logger, BrightMagenta))
	case "panic":
		log.Panic(FormatLogger(logger, BrightRed))
	default:
		log.Print(FormatLogger(logger, BrightCyan))
	}
}

// Simple outputs a log with the default information.
func Simple() {
	// Disable default log timestamp
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Get information about the file and the line
	_, filename, line, _ := runtime.Caller(1)

	logger := &LoggerContent{
		Filename:      filename,
		Line:          strconv.Itoa(line),
		Timestamp:     time.Now().Format(time.RFC3339),
		LoggerStatus:  &LoggerStatus{},
		LoggerMessage: &LoggerMessage{},
		LoggerFault:   &LoggerFault{},
	}

	err := SaveLogger(logger)
	if err != nil {
		log.Fatalf("Unable to save logger to log file: %v", err)
	}

	log.Print(FormatLogger(logger, BrightCyan))
}

// Status outputs a log with a given status.
func Status(status string) {
	// Disable default log timestamp
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Get information about the file and the line
	_, filename, line, _ := runtime.Caller(1)

	logger := &LoggerContent{
		Filename:      filename,
		Line:          strconv.Itoa(line),
		Timestamp:     time.Now().Format(time.RFC3339),
		LoggerStatus:  &LoggerStatus{Status: status},
		LoggerMessage: &LoggerMessage{},
		LoggerFault:   &LoggerFault{},
	}

	err := SaveLogger(logger)
	if err != nil {
		log.Fatalf("Unable to save logger to log file: %v", err)
	}

	switch logger.LoggerStatus.Status {
	case "success":
		log.Print(FormatLogger(logger, BrightGreen))
	case "info":
		log.Print(FormatLogger(logger, BrightBlue))
	case "error":
		log.Print(FormatLogger(logger, BrightYellow))
	case "fatal":
		log.Fatal(FormatLogger(logger, BrightMagenta))
	case "panic":
		log.Panic(FormatLogger(logger, BrightRed))
	default:
		log.Print(FormatLogger(logger, BrightCyan))
	}
}

// Message outputs a log with a given message.
func Message(message string) {
	// Disable default log timestamp
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Get information about the file and the line
	_, filename, line, _ := runtime.Caller(1)

	logger := &LoggerContent{
		Filename:      filename,
		Line:          strconv.Itoa(line),
		Timestamp:     time.Now().Format(time.RFC3339),
		LoggerStatus:  &LoggerStatus{},
		LoggerMessage: &LoggerMessage{Message: message},
		LoggerFault:   &LoggerFault{},
	}

	err := SaveLogger(logger)
	if err != nil {
		log.Fatalf("Unable to save logger to log file: %v", err)
	}

	log.Print(FormatLogger(logger, BrightCyan))
}

// Error outputs a log with a given fault.
func Error(fault error) {
	// Disable default log timestamp
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Get information about the file and the line
	_, filename, line, _ := runtime.Caller(1)

	logger := &LoggerContent{
		Filename:      filename,
		Line:          strconv.Itoa(line),
		Timestamp:     time.Now().Format(time.RFC3339),
		LoggerStatus:  &LoggerStatus{},
		LoggerMessage: &LoggerMessage{},
		LoggerFault:   &LoggerFault{Fault: fmt.Sprint(fault)},
	}

	err := SaveLogger(logger)
	if err != nil {
		log.Fatalf("Unable to save logger to log file: %v", err)
	}

	log.Print(logger.LoggerFault.Fault)
	log.Print(FormatLogger(logger, BrightYellow))
}

// Fatal outputs a log with a given fault.
func Fatal(fault error) {
	// Disable default log timestamp
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Get information about the file and the line
	_, filename, line, _ := runtime.Caller(1)

	logger := &LoggerContent{
		Filename:      filename,
		Line:          strconv.Itoa(line),
		Timestamp:     time.Now().Format(time.RFC3339),
		LoggerStatus:  &LoggerStatus{},
		LoggerMessage: &LoggerMessage{},
		LoggerFault:   &LoggerFault{Fault: fmt.Sprint(fault)},
	}

	err := SaveLogger(logger)
	if err != nil {
		log.Fatalf("Unable to save logger to log file: %v", err)
	}

	log.Print(FormatLogger(logger, BrightMagenta))
}

// Panic outputs a log with a given fault.
func Panic(fault error) {
	// Disable default log timestamp
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Get information about the file and the line
	_, filename, line, _ := runtime.Caller(1)

	logger := &LoggerContent{
		Filename:      filename,
		Line:          strconv.Itoa(line),
		Timestamp:     time.Now().Format(time.RFC3339),
		LoggerStatus:  &LoggerStatus{},
		LoggerMessage: &LoggerMessage{},
		LoggerFault:   &LoggerFault{Fault: fmt.Sprint(fault)},
	}

	err := SaveLogger(logger)
	if err != nil {
		log.Fatalf("Unable to save logger to log file: %v", err)
	}

	log.Print(FormatLogger(logger, BrightRed))
}
