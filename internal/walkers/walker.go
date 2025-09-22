package walkers

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Talos-hub/ZibraGo/internal/apperrors"
)

type Walker struct {
	pathToDir string
}

func NewWalker(pathTodir string) *Walker {
	return &Walker{pathToDir: pathTodir}
}

func (w *Walker) Walk() ([]string, error) {
	info, err := os.Stat(w.pathToDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, apperrors.NewAppError("path not exist", "Walk", apperrors.E_OPEN, err)
		}
		return nil, apperrors.NewAppError("error get info", "Walk", apperrors.E_OPEN, err)
	}

	if !info.IsDir() {
		return nil, apperrors.NewAppError("path is not dir", "Walk", apperrors.E_OPEN, err)
	}

	pathes := make([]string, 0, 100)

	err = filepath.Walk(w.pathToDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			pathes = append(pathes, path)
		}
		return nil
	})

	if err != nil {
		return nil, apperrors.NewAppError("error walking", "Walk", apperrors.E_READ, err)
	}

	return pathes, nil

}
