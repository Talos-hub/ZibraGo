package configuration

import (
	"encoding/json"
	"os"

	"github.com/Talos-hub/ZibraGo/internal/apperrors"
)

const name = "zipdir.json"
const DEFAULT_PATH = "."

type Config struct {
	ZipDir string `json:"zipdir"`
}

func GetDir() (*Config, error) {
	file, err := os.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, apperrors.NewAppError("error open config file", "GetDir", apperrors.E_OPEN, err)
	}
	defer file.Close()

	config := &Config{}

	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, apperrors.NewAppError("error decoding", "GetDir", apperrors.E_PARSE, err)
	}
	return config, nil
}

func NewPath(path string) error {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return apperrors.NewAppError("error open or create a config file", "NewPath", apperrors.E_OPEN, err)
	}

	defer file.Close()

	config := &Config{ZipDir: path}

	err = json.NewEncoder(file).Encode(config)
	if err != nil {
		return apperrors.NewAppError("error encoding", "NewPath", apperrors.E_CONF, err)
	}

	return nil
}
