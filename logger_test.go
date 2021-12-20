package teavet

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

// TestStartLog checks whether or not the log file has been created
func TestStartLog(t *testing.T) {
	err := StartLog()
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

	t.Log("\n✅ StartLog test passed...\n")
}

// TestFindLog checks whether or not the log file is found
func TestFindLog(t *testing.T) {
	logFile, err := FindLog()
	if err != nil {
		t.Fatalf("\n❌ FindLog failed: %v\n", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("\n❌ Unable to get working directory: %v\n", err)
	}

	if logFile != path.Join(dir, "teavet.log") {
		t.Fatalf("\n❌ Unable to match log file: %v\n", err)
	}

	t.Log("\n✅ FindLog test passed...\n")
}

// TestBuildLog checks whether or not the log file has content
func TestBuildLog(t *testing.T) {
	logContent, err := BuildLog(&Optional{Status: "success", Message: "testing", Error: nil})
	if err != nil {
		t.Fatalf("\n❌ BuildLog failed: %v\n", err)
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

	if content["filename"] != logContent.Filename {
		t.Fatalf("\n❌ Unable to match log filename: %v\n", err)
	}

	t.Log("\n✅ BuildLog test passed...\n")
}

// TestLogger returns the output of Logger
func TestLogger(t *testing.T) {
	Logger("success", "testing logger", nil)

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("\n❌ Unable to get working directory: %v\n", err)
	}

	// Remove previously created log file
	err = os.Remove(path.Join(dir, "teavet.log"))
	if err != nil {
		t.Fatalf("\n❌ Unable to remove log file: %v\n", err)
	}
}
