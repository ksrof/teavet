package teavet

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestStartLogger checks whether or not the log file has been created
func TestStartLogger(t *testing.T) {
	err := StartLogger()
	if err != nil {
		t.Fatalf("\n❌ StartLog failed: %v\n", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("\n❌ Unable to get working directory: %v\n", err)
	}

	// Look for the log file
	file, err := filepath.Glob(path.Join(dir, "teavet.log"))
	if err != nil {
		t.Fatalf("\n❌ Unable to find the log file: %v\n", err)
	}

	filename := strings.Join(file, "")

	if filename != path.Join(dir, "teavet.log") {
		t.Fatalf("\n❌ Unable to match log file: %v\n", err)
	}

	t.Log("\n✅ StartLogger test passed...\n")
}

// TestFindLogger checks whether or not the log file is found
func TestFindLogger(t *testing.T) {
	loggerFile, err := FindLogger()
	if err != nil {
		t.Fatalf("\n❌ FindLog failed: %v\n", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("\n❌ Unable to get working directory: %v\n", err)
	}

	if loggerFile != path.Join(dir, "teavet.log") {
		t.Fatalf("\n❌ Unable to match log file: %v\n", err)
	}

	t.Log("\n✅ FindLog test passed...\n")
}

// TestSaveLogger checks whether or not the log file has content
func TestSaveLogger(t *testing.T) {
	// Get information about the file and the line
	_, filename, line, _ := runtime.Caller(1)

	logger := &LoggerContent{
		Filename:      filename,
		Line:          strconv.Itoa(line),
		Timestamp:     time.Now().Format(time.RFC3339),
		LoggerStatus:  &LoggerStatus{Status: "Success"},
		LoggerMessage: &LoggerMessage{Message: "Teavet is up and running!"},
		LoggerFault:   &LoggerFault{},
	}

	err := SaveLogger(logger)
	if err != nil {
		t.Fatalf("\n❌ Unable to save logger to log file: %v", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("\n❌ Unable to get working directory: %v\n", err)
	}

	fileContent, err := os.ReadFile(path.Join(dir, "teavet.log"))
	if err != nil {
		t.Fatalf("\n❌ Unable to read log file: %v\n", err)
	}

	var content map[string]interface{}
	err = json.Unmarshal(fileContent, &content)
	if err != nil {
		t.Fatalf("\n❌ Unable to unmarshal log file content: %v\n", err)
	}

	if content["filename"] != logger.Filename {
		t.Fatalf("\n❌ Unable to match log filename: %v\n", err)
	}

	t.Log("\n✅ BuildLog test passed...\n")
}

// TestLogger returns the output of Complete logger method
func TestLogger(t *testing.T) {
	// Complete log type
	Complete("success", "Teavet is up and running!", nil)

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("\n❌ Unable to get working directory: %v\n", err)
	}

	// Remove previously created log file
	err = os.Remove(path.Join(dir, "teavet.log"))
	if err != nil {
		t.Fatalf("\n❌ Unable to remove log file: %v\n", err)
	}

	t.Log("\n✅ TestLogger test passed...\n")
}
