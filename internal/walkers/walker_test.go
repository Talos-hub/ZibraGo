package walkers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWalker_Pathes(t *testing.T) {
	tempDir := t.TempDir()

	files := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		files = append(files, filepath.Join(tempDir, fmt.Sprintf("test%d.txt", i)))
	}

	for _, file := range files {
		err := os.WriteFile(file, []byte("Test data"), 0644)
		if err != nil {
			t.Fatalf("failed create test files: %v", err)
		}

	}

	wakler := NewWalker(tempDir)

	data, err := wakler.Walk()
	if err != nil {
		t.Errorf("failed walk: %v", err)
	}

	found := make(map[string]bool, 100)

	for _, file := range data {
		found[file] = true
	}

	for _, expected := range files {
		if !found[expected] {
			t.Errorf("Expected: %s", expected)
		}
	}
}

func BenchmarkWalker(b *testing.B) {
	// Create test structure once outside the benchmark loop
	tempDir := b.TempDir()
	createTestFileStructure(b, tempDir, 10, 10) // 10 dirs, 10 files each

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		walker := NewWalker(tempDir)
		_, err := walker.Walk()
		if err != nil {
			b.Fatalf("Walk failed: %v", err)
		}
	}
}

func TestWalkerMatch_slice(t *testing.T) {
	pathes := []string{
		"test.exe", "test.txt",
		"test.json", "test.dll", "test.png",
		"test.jpg", "test.jepg",
	}
	extentions := map[string]bool{
		".jpg":   true,
		".png":   true,
		".jepg:": true,
	}

	//act
	slice := match(pathes, extentions)

	//assert

	for _, item := range slice {
		file := filepath.Base(item)
		n := strings.Index(file, ".")

		ex := file[n:]

		if _, ok := extentions[ex]; ok {
			t.Fatalf("expected there is not a file with the extension: %s", ex)
		}
	}

}

func BenchmarkMatch(b *testing.B) {
	// arange
	pathes := []string{
		"test.exe", "test.txt",
		"test.json", "test.dll",
		"test.png", "test2.jpg",
		"test.jpg", "test.jepg",
		"test3.jpg", "test4.jpg",
		"test5.jpg", "test6.jpg",
		"test6.jpg", "test8.jpg",
		"test9.jpg", "test10.jpg",
	}
	extentions := map[string]bool{
		".jpg":   true,
		".png":   true,
		".jepg:": true,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = match(pathes, extentions)
	}
}

// Helper function to create realistic test structure
func createTestFileStructure(b *testing.B, baseDir string, numDirs, filesPerDir int) {
	b.Helper()

	// Create some files in root directory
	for i := 0; i < filesPerDir; i++ {
		filePath := filepath.Join(baseDir, fmt.Sprintf("rootfile%d.txt", i))
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			b.Fatalf("failed to create root file: %v", err)
		}
	}

	// Create subdirectories with files
	for i := 0; i < numDirs; i++ {
		subdirPath := filepath.Join(baseDir, fmt.Sprintf("subdir%d", i))
		err := os.MkdirAll(subdirPath, 0755)
		if err != nil {
			b.Fatalf("failed to create subdirectory: %v", err)
		}

		// Create files in this subdirectory
		for j := 0; j < filesPerDir; j++ {
			filePath := filepath.Join(subdirPath, fmt.Sprintf("file%d.txt", j))
			err := os.WriteFile(filePath, []byte("test content"), 0644)
			if err != nil {
				b.Fatalf("failed to create file: %v", err)
			}
		}

		// Create nested subdirectory (more realistic)
		if i%3 == 0 { // Every 3rd directory has a nested subdir
			nestedPath := filepath.Join(subdirPath, "nested")
			err := os.MkdirAll(nestedPath, 0755)
			if err != nil {
				b.Fatalf("failed to create nested directory: %v", err)
			}

			for k := 0; k < filesPerDir/2; k++ {
				filePath := filepath.Join(nestedPath, fmt.Sprintf("nested%d.txt", k))
				err := os.WriteFile(filePath, []byte("test content"), 0644)
				if err != nil {
					b.Fatalf("failed to create nested file: %v", err)
				}
			}
		}
	}
}
