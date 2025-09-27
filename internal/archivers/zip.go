package archivers

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Talos-hub/ZibraGo/internal/apperrors"
)

type ZipArchiver struct {
	pathToFolder string
	basePath     string // Base path for relative calculations (folder)
	nameZip      string
	maxWorkers   int
	mu           *sync.Mutex
	wg           *sync.WaitGroup
}

// NewZipArchiver returns a pointer to ZipArchiver
func NewZipArchiver(pathTofolder string, basePath string, nameZip string, maxWorkers int) *ZipArchiver {
	return &ZipArchiver{
		pathToFolder: pathTofolder,
		basePath:     basePath,
		nameZip:      nameZip,
		maxWorkers:   maxWorkers,
		wg:           &sync.WaitGroup{},
		mu:           &sync.Mutex{},
	}
}

// Start create new archive and add files
func (a *ZipArchiver) Start(paths ...string) (string, error) {
	// Check that folder path is correct
	err := a.isFolder()
	if err != nil {
		return "", apperrors.NewAppError("error to start zip Archiver", "isFolder", apperrors.E_OPEN, err)
	}

	path := filepath.Join(a.pathToFolder, a.nameZip)

	// Create new Zip Archive
	zipfile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return "", apperrors.NewAppError("error to start zip Archiver, the file already exist", "Start", apperrors.E_CREATE, err)
		}
		return "", apperrors.NewAppError("error open a file", "Start", apperrors.E_OPEN, err)
	}

	defer zipfile.Close()

	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()

	// add single file to zip
	if len(paths) == 1 {
		err := addFileToZip(zipWriter, paths[0], a.basePath, a.mu)
		if err != nil {

			os.Remove(zipfile.Name())
			return "", apperrors.NewAppError("error add a single file to zip", "addFileToZip", apperrors.E_ADD, err)
		}
		return zipfile.Name(), nil
	}
	// jobs
	fileChan := make(chan string, len(paths))
	errChan := make(chan error, len(paths))

	// start workers
	for i := 0; i < a.maxWorkers; i++ {
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()

			for filepath := range fileChan {
				if err := addFileToZip(zipWriter, filepath, a.basePath, a.mu); err != nil {
					errChan <- apperrors.NewAppError("error add files to zip", "addFileToZip", apperrors.E_ADD, err)
				}
			}
		}()
	}

	// send files to chan
	for _, filePath := range paths {
		fileChan <- filePath
	}
	// signal that there is no more files
	close(fileChan)
	a.wg.Wait()

	close(errChan)
	sliceErrors := make([]error, 0, len(errChan))

	for err := range errChan {
		if err != nil {
			sliceErrors = append(sliceErrors, err)
		}
	}

	// if there is an error, so is deletes the zip file
	if len(sliceErrors) > 0 {
		os.Remove(zipfile.Name())
		return "", fmt.Errorf("multiple errors occurred: %v", sliceErrors)
	}

	return zipfile.Name(), nil
}

// addFileToZip add a single file to zip
func addFileToZip(zipWriter *zip.Writer, filename string, basepath string, mu *sync.Mutex) error {
	// open a file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error open a file: %w", err)
	}
	defer file.Close()

	// get file info to set zip header
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error get file info from: %s, %w", filename, err)
	}

	// set headerinfo
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("error to set file info header: %w", err)
	}

	relativePath, err := filepath.Rel(basepath, filename)
	if err != nil {
		return fmt.Errorf("error calculating relative path for: %s, %w", filename, err)
	}

	relativePath = filepath.ToSlash(relativePath)

	header.Method = zip.Deflate
	header.Name = relativePath
	// lock function for safe zipWriter
	mu.Lock()
	defer mu.Unlock()

	// create header
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("error create zip  header for: %s, %w", filename, err)
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return fmt.Errorf("error add a file %s to zip: %w", filename, err)
	}
	return nil
}

// isFolder is helper function, check that path of folder exist
func (z *ZipArchiver) isFolder() error {
	info, err := os.Stat(z.pathToFolder)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path is not exist")
		}
		return fmt.Errorf("error get info, path: %s, err: %w", z.pathToFolder, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("it is not a folder, path: %s, err: %w", z.pathToFolder, err)
	}

	return nil
}

func IsZip(file string) bool {

	if len(file) < 4 {
		return false
	}

	return strings.Contains(file, ".zip")
}
