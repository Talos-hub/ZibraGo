package archivers

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func createTestFiles(t *testing.T, dir string, files []string) {
	t.Helper()

	for _, fileName := range files {
		path := filepath.Join(dir, fileName)

		// Create directory if needed
		dirPath := filepath.Dir(path)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Error creating directory: %v", err)
		}

		err := os.WriteFile(path, []byte("test data for "+fileName), 0644)
		if err != nil {
			t.Fatalf("Error create test files: %v", err)
		}
	}
}

func verifyZipContents(t *testing.T, zipPath string, expected []string) {
	t.Helper()

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("Failed to open zip: %v", err)
	}
	defer reader.Close()

	foundFiles := make(map[string]bool)
	for _, file := range reader.File {
		foundFiles[file.Name] = true
		t.Logf("Found in zip: %s", file.Name) // Debug logging
	}

	for _, expectedFile := range expected {
		if !foundFiles[expectedFile] {
			t.Errorf("Expected file %s not found in zip. Found: %v", expectedFile, getKeys(foundFiles))
		}
	}
}

// Helper function to get map keys for error message
func getKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TestNewZipArchiver(t *testing.T) {
	temp := t.TempDir()

	// Fixed parameter order: pathToFolder, basePath, nameZip, maxWorkers
	archiver := NewZipArchiver(temp, temp, "test.zip", 4)

	if archiver == nil {
		t.Error("Expected non-nil archiver")
		return
	}

	if archiver.maxWorkers != 4 {
		t.Errorf("Expected 4 workers, got: %d", archiver.maxWorkers)
	}

	if archiver.nameZip != "test.zip" {
		t.Errorf("Expected name test.zip got: %s", archiver.nameZip)
	}
	if archiver.pathToFolder != temp {
		t.Errorf("Expected path to folder: %s, got: %s", temp, archiver.pathToFolder)
	}
}

func TestZipArchiver_SingleFile(t *testing.T) {
	tempDir := t.TempDir()

	files := []string{"test.txt"}
	createTestFiles(t, tempDir, files)

	// Fixed parameter order
	archiver := NewZipArchiver(tempDir, tempDir, "single.zip", 2)

	file := filepath.Join(tempDir, files[0])
	zipPath, err := archiver.Start(file)
	if err != nil {
		t.Fatalf("Failed to create archive: %v", err)
	}

	// Verify zip file exists and has content
	info, err := os.Stat(zipPath)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	if info.Size() == 0 {
		t.Error("Zip should not be empty")
	}

	// Verify the zip contains our file
	verifyZipContents(t, zipPath, files)
}

func TestZipArchiver_MultipleFiles(t *testing.T) {
	tempDir := t.TempDir()

	testFiles := []string{"test1.txt", "test2.txt", "test3.txt", "test4.txt"}
	createTestFiles(t, tempDir, testFiles)

	// Fixed parameter order
	archiver := NewZipArchiver(tempDir, tempDir, "multiple.zip", 2)

	files := make([]string, 0, len(testFiles))
	for _, path := range testFiles {
		files = append(files, filepath.Join(tempDir, path))
	}

	zipPath, err := archiver.Start(files...)
	if err != nil {
		t.Fatalf("Failed to create archive: %v", err)
	}

	verifyZipContents(t, zipPath, testFiles)
}

func TestZipArchiver_InvalidFolder(t *testing.T) {
	// Fixed parameter order
	archiver := NewZipArchiver("xxx/xxx/path", ".", "test.zip", 2)

	_, err := archiver.Start()
	if err == nil {
		t.Error("Expected error for invalid folder")
	}
}

func TestZipArchiver_NonExistentFile(t *testing.T) {
	tempDir := t.TempDir()

	// Fixed parameter order
	archiver := NewZipArchiver(tempDir, tempDir, "test.zip", 2)

	_, err := archiver.Start("/non/existentfile/file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestZipArchiver_EmptyFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Test with empty file list
	archiver := NewZipArchiver(tempDir, tempDir, "empty.zip", 2)

	zipPath, err := archiver.Start()
	if err != nil {
		t.Fatalf("Should handle empty file list: %v", err)
	}

	// Verify empty zip was created
	info, err := os.Stat(zipPath)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	if info.Size() == 0 {
		t.Log("Empty zip created (expected for no files)")
	}
}

func TestZipArchiver_WithSubdirectories(t *testing.T) {
	tempDir := t.TempDir()

	// Create files in subdirectories
	files := []string{
		"file1.txt",
		"subdir/file2.txt",
		"subdir/nested/file3.txt",
	}
	createTestFiles(t, tempDir, files)

	archiver := NewZipArchiver(tempDir, tempDir, "subdirs.zip", 2)

	fullPaths := make([]string, 0, len(files))
	for _, file := range files {
		fullPaths = append(fullPaths, filepath.Join(tempDir, file))
	}

	zipPath, err := archiver.Start(fullPaths...)
	if err != nil {
		t.Fatalf("Failed to create archive with subdirectories: %v", err)
	}

	verifyZipContents(t, zipPath, files)
}

// Fixed benchmark with proper setup
func BenchmarkZipArchiverStart(b *testing.B) {
	tempDir := b.TempDir()

	// Create benchmark files
	var testFiles []string
	for i := 0; i < 100; i++ {
		testFiles = append(testFiles, fmt.Sprintf("benchfile%d.txt", i))
	}

	// Create files once before benchmarking
	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		// Create smaller file for faster benchmarking (100KB instead of 100MB)
		content := make([]byte, 1024*100)
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			b.Fatalf("Failed to create benchmark file: %v", err)
		}
	}

	var filePaths []string
	for _, file := range testFiles {
		filePaths = append(filePaths, filepath.Join(tempDir, file))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		archiver := NewZipArchiver(tempDir, tempDir, fmt.Sprintf("bench%d.zip", i), 4)
		_, err := archiver.Start(filePaths...)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}
