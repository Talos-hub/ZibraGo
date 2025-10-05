package walkers

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Talos-hub/ZibraGo/internal/apperrors"
	"github.com/Talos-hub/ZibraGo/internal/configuration"
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
		return nil, apperrors.NewAppError("path is not dir", "Walk", apperrors.E_OPEN, nil)
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

	if len(pathes) == 0 {
		return pathes, nil
	}

	m, err := configuration.GetExtentions()
	if err != nil {
		if os.IsNotExist(err) {
			return pathes, nil
		}
		return nil, apperrors.NewAppError("error get extentions", "Walk", apperrors.E_CONF, err)
	}

	// matching
	pathes = match(pathes, m)
	return pathes, nil

}

// match match paths and extensions if an element from the path is
// to an element from the extension
// the element will not be added to the slice
func match(paths []string, extentions map[string]bool) []string {
	slice := make([]string, 0, len(paths)) // determine a slice

	if len(extentions) == 0 {
		return paths
	}

	for _, item := range paths {
		file := filepath.Base(item)

		n := strings.Index(file, ".")
		// cut unuseful
		extention := file[n:]
		// match
		if _, ok := extentions[extention]; ok {
			continue
		}

		slice = append(slice, item)
	}

	return slice
}
