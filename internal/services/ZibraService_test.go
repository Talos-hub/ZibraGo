package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type WalkerMock struct {
	mock.Mock
}

func (w *WalkerMock) Walk() ([]string, error) {
	args := w.Called()
	return args.Get(0).([]string), args.Error(1)
}

type ArchiverMock struct {
	mock.Mock
}

func (a *ArchiverMock) Start(paths ...string) (string, error) {
	args := a.Called(paths)
	return args.String(0), args.Error(1)
}

type ApiMock struct {
	mock.Mock
}

func (a *ApiMock) Check() error {
	args := a.Called()
	return args.Error(0)
}

func (a *ApiMock) UpLoadFile(file *os.File) error {
	args := a.Called(file)
	return args.Error(0)
}

type LoggerMock struct {
	mock.Mock
}

func (l *LoggerMock) Info(msg string, arg ...any) {
	allArgs := append([]interface{}{msg}, arg...)
	l.Called(allArgs...)
}

func (l *LoggerMock) Warn(msg string, arg ...any) {
	allArgs := append([]interface{}{msg}, arg...)
	l.Called(allArgs...)
}

func (l *LoggerMock) Debug(msg string, arg ...any) {
	allArgs := append([]interface{}{msg}, arg...)
	l.Called(allArgs...)
}

func (l *LoggerMock) Error(msg string, arg ...any) {
	allArgs := append([]interface{}{msg}, arg...)
	l.Called(allArgs...)
}

func TestZibraServiceRun_Run_Success(t *testing.T) {
	archiver := &ArchiverMock{}
	walker := &WalkerMock{}
	api := &ApiMock{}
	logger := &LoggerMock{}

	dir := t.TempDir()

	testfile1 := filepath.Join(dir, "test1.txt")
	testfile2 := filepath.Join(dir, "test2.txt")
	os.WriteFile(testfile1, []byte("test content 1"), 0644)
	os.WriteFile(testfile2, []byte("test content 2"), 0644)

	expectedFiles := []string{testfile1, testfile2}
	archivePath := filepath.Join(dir, "archive.zip")

	file, err := os.OpenFile(archivePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal("Failed to create test zip")
	}
	file.Close()

	// set mock
	walker.On("Walk").Return(expectedFiles, nil)
	archiver.On("Start", expectedFiles).Return(archivePath, nil)
	api.On("Check").Return(nil)

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Maybe()
	logger.On("Warn", mock.Anything, mock.Anything, mock.Anything).Maybe()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Maybe()
	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything).Maybe()

	api.On("UpLoadFile", mock.MatchedBy(func(file *os.File) bool {
		return file.Name() == archivePath
	})).Return(nil)

	zibra, err := NewZibra(walker, archiver, api, logger)
	if err != nil {
		t.Errorf("Failed to create new zibra service, expected nil err: %v", err)
	}

	if zibra == nil {
		t.Fatalf("Expected non-nil zibra service: %v", zibra)
	}

	err = zibra.Run()

	assert.NoError(t, err)

	api.AssertCalled(t, "Check")
	walker.AssertCalled(t, "Walk")
	archiver.AssertCalled(t, "Start", expectedFiles)
	api.AssertCalled(t, "UpLoadFile", mock.MatchedBy(func(file *os.File) bool {
		return file.Name() == archivePath
	}))
}
