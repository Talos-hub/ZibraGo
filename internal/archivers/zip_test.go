package archivers

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func crateTestFiles(t *testing.T, dir string, files []string) {
	t.Helper()

	for _, fileName := range files {
		path := filepath.Join(dir, fileName)

		err := os.WriteFile(path, []byte("test data for "+fileName), 0644)
		if err != nil {
			t.Fatalf("Error create test files: %v", err)
		}
	}
}

func verifyZipContents(t *testing.T, ziPath string, exptected []string) {
	t.Helper()

	reader, err := zip.OpenReader(ziPath)
	if err != nil {
		t.Errorf("Failed to open zip: %v", err)
	}
	defer reader.Close()

	foundFiles := make(map[string]bool)

	for _, file := range reader.File {
		foundFiles[file.Name] = true
	}

	for _, file := range exptected {
		if !foundFiles[filepath.Base(file)] {
			t.Errorf("Expected file %s, not found", file)
		}
	}
}

func TestNewZipArchivers(t *testing.T) {
	temp := t.TempDir()

	archivers := NewZipArchiver(temp, "test.zip", 4)

	if archivers == nil {
		t.Error("Expected non-nil archiver")
		return
	}

	if archivers.maxWorkers != 4 {
		t.Errorf("Expdected 4 wrokers, got: %d", archivers.maxWorkers)
	}

	if archivers.nameZip != "test.zip" {
		t.Errorf("Expected name test.zip got: %s", archivers.nameZip)
	}
	if archivers.pathToFolder != temp {
		t.Errorf("Expected path to folder: %s, got: %s", temp, archivers.pathToFolder)
	}
}

func TestZipArchiver_SingleFile(t *testing.T) {
	tempDir := t.TempDir()

	files := []string{"test.txt"}

	crateTestFiles(t, tempDir, files)

	archivers := NewZipArchiver(tempDir, "single.zip", 2)

	file := filepath.Join(tempDir, files[0])

	zipPath, err := archivers.Start(file)
	if err != nil {
		t.Errorf("Failed to create archive, %v", err)
	}

	zip, err := os.Open(zipPath)
	if err != nil {
		t.Errorf("Failed to open zip file, %v", err)
	}
	defer zip.Close()

	info, err := zip.Stat()
	if err != nil {
		t.Errorf("Failed to get file info: %v", err)
	}

	if info.Size() == 0 {
		t.Error("Zip should be not empty")
	}

	if _, err := os.Stat(filepath.Join(tempDir, "single.zip")); os.IsNotExist(err) {
		t.Errorf("Zip file was not created at the disk, %v", err)
	}
}

func TestZipArchiver_MultipleFiles(t *testing.T) {
	tempDir := t.TempDir()

	testFiles := []string{"test1.txt", "text2.txt", "text3.txt", "text4.txt"}

	archiver := NewZipArchiver(tempDir, "multiple.zip", 2)

	files := make([]string, 0, len(testFiles))

	for _, path := range testFiles {
		files = append(files, filepath.Join(tempDir, path))
	}

	crateTestFiles(t, tempDir, testFiles)

	zipPath, err := archiver.Start(files...)
	if err != nil {
		t.Errorf("failed to create archive: %v", err)
	}

	zip, err := os.Open(zipPath)
	if err != nil {
		t.Errorf("Failed to open zip file, %v", err)
	}
	defer zip.Close()

	verifyZipContents(t, zip.Name(), testFiles)
}

func TestZipArchiver_InvalidFolder(t *testing.T) {
	zip := NewZipArchiver("xxx/xxx/path", "test.zip", 2)

	_, err := zip.Start()
	if err == nil {
		t.Errorf("Expected err, got: %v", err)
	}
}

func TestZipArchiver_NonExistenFile(t *testing.T) {
	tempDir := t.TempDir()

	zip := NewZipArchiver(tempDir, "test.zip", 2)

	_, err := zip.Start("/non/existenfile/file.txt")

	if err == nil {
		t.Errorf("Expected err, got: %v", err)
	}
}

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
		// Create 1MB file
		content := make([]byte, 1024*1024*100)
		os.WriteFile(filePath, content, 0644)
	}

	var filePaths []string
	for _, file := range testFiles {
		filePaths = append(filePaths, filepath.Join(tempDir, file))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		archiver := NewZipArchiver(tempDir, fmt.Sprintf("bench%d.zip", i), 4)
		zipFile, err := archiver.Start(filePaths...)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
		fmt.Println(zipFile)
	}

}
