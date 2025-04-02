package common

import (
	// "encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	// "text/template"
)

// TestGenerateConfig tests the GenerateConfig function
func TestGenerateConfig(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "test-config")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after the test

	// Create a dummy template file
	templateFile := filepath.Join(tempDir, "/test_template.cfg")
	err = ioutil.WriteFile(templateFile, []byte("Test template content with data: {{.TestData}}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Create a dummy output directory
	outputDir := filepath.Join(tempDir, "output")
	err = os.Mkdir(outputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	outputFile := filepath.Join(outputDir, "/test_output.cfg")

	// Test data
	data := map[string]string{"TestData": "Hello, Test!"}

	// Execute the function
	err = GenerateConfig(templateFile, outputFile, data)
	if err != nil {
		t.Fatalf("GenerateConfig failed: %v", err)
	}

	// Verify that the output file was created
	_, err = os.Stat(outputFile)
	if os.IsNotExist(err) {
		t.Errorf("Output file was not created")
	}

	// Read the output file and verify its content
	outputData, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedOutput := "Test template content with data: Hello, Test!"
	if string(outputData) != expectedOutput {
		t.Errorf("Output file content is incorrect. Expected: '%s', Got: '%s'", expectedOutput, string(outputData))
	}
}

func TestGenerateConfig_TemplateNotFound(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "test-config")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after the test

	outputFile := filepath.Join(tempDir, "test_output.cfg")
	data := map[string]string{"TestData": "Hello, Test!"}

	// Execute the function with a non-existent template file
	err = GenerateConfig("non_existent_template.cfg", filepath.Base(outputFile), data)
	if err == nil {
		t.Errorf("GenerateConfig should have failed with a non-existent template")
	}
}

func TestGenerateConfig_OutputCreationError(t *testing.T) {
    // Create a temporary directory for test files
    tempDir, err := os.MkdirTemp("", "test-config")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir) // Clean up after the test

    // Create a dummy template file
    templateFile := filepath.Join(tempDir, "test_template.cfg")
    err = ioutil.WriteFile(templateFile, []byte("Test template content with data: {{.TestData}}"), 0644)
    if err != nil {
        t.Fatalf("Failed to create template file: %v", err)
    }

    // Create a dummy output directory that cannot be created
    outputFile := "/nonexistent/test_output.cfg" // This should cause an error

    // Test data
    data := map[string]string{"TestData": "Hello, Test!"}

    // Execute the function
    err = GenerateConfig(filepath.Base(templateFile), filepath.Base(outputFile), data)
    if err == nil {
        t.Errorf("GenerateConfig should have failed due to output file creation error")
    }
}

func TestGenerateConfig_TemplateParseError(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "test-config")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after the test

	// Create a template file with an error
	templateFile := filepath.Join(tempDir, "error_template.cfg")
	err = ioutil.WriteFile(templateFile, []byte("Test template content with bad syntax: {{.BadData"), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	outputFile := filepath.Join(tempDir, "test_output.cfg")
	data := map[string]string{"TestData": "Hello, Test!"}

	// Execute the function
	err = GenerateConfig(filepath.Base(templateFile), filepath.Base(outputFile), data)
	if err == nil {
		t.Errorf("GenerateConfig should have failed due to template parse error")
	}
}

func TestDebugData(t *testing.T) {
	// This test mainly checks that DebugData doesn't panic.
	// You could add more sophisticated checks if needed,
	// like capturing the log output and verifying its content.

	data := map[string]string{"key": "value"}
	DebugData(data) // Call the function

	// If the function doesn't panic, the test passes.
	// You might want to add assertions about the log output if you need to be more specific.
}