package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Talos-hub/ZibraGo/internal/apperrors"
)

const name = "zipdir.json"
const DEFAULT_PATH = "."
const extention = "extentions.json"

type Config struct {
	ZipDir string `json:"zipdir"`
}

func GetDir() (*Config, error) {
	file, err := os.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, apperrors.NewAppError("error open configuration file", "GetDir", apperrors.E_OPEN, err)
	}
	defer file.Close()

	config := &Config{}

	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, apperrors.NewAppError("error decoding", "GetDir", apperrors.E_PARSE, err)
	}
	return config, nil
}

func GetExtentions() (map[string]bool, error) {
	file, err := os.Open(extention)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, apperrors.NewAppError("error open a file", "GetExtentions", apperrors.E_OPEN, err)
	}
	defer file.Close()

	m := make(map[string]bool)

	err = json.NewDecoder(file).Decode(&m)
	if err != nil {
		return nil, apperrors.NewAppError("error decoding", "GetExtentions", apperrors.E_PARSE, err)
	}

	return m, nil

}

// CreateExtentios reads user input and add extentions to a file
func CreateExtentions() error {
	file, err := os.OpenFile(extention, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return apperrors.NewAppError("error create or open a file", "CreateExtentios", apperrors.E_CREATE, err)
	}
	defer file.Close()

	var input string

	fmt.Println("Please enter extentions like: .exe .txt .json")
	_, err = fmt.Scan(&input)
	if err != nil {
		return apperrors.NewAppError("error scanning user input", "CreateExtentios", apperrors.E_PARSE, err)
	}

	m, err := cutExtentions(input)
	if err != nil {
		return apperrors.NewAppError("error create extentions, wrong format", "CreateExtentions", apperrors.E_CREATE, err)
	}

	err = json.NewEncoder(file).Encode(m)
	if err != nil {
		return apperrors.NewAppError("error encoding extentions", "CreatingExtentions", apperrors.E_CREATE, err)
	}

	return nil

}

// cutExtention cut extentions and validate them
func cutExtentions(ex string) (map[string]bool, error) {
	if len(ex) == 0 {
		return nil, errors.New("error, user input is empty")
	}
	if strings.Contains(ex, ",") {
		return nil, errors.New("error format, you should write like: .json .exe .txt INSTED .json,.exe,.txt, or .json, .txt, .exe")
	}
	items := strings.Split(ex, " ")
	m := make(map[string]bool, len(items))

	for _, s := range items {
		n := strings.Index(s, ".")
		// if there is no dot, so it is not an extention
		if n == -1 {
			return nil, fmt.Errorf("error, this: %s is not an extention", s)
		}

		// cut extention
		extention := s[n:]

		if len(extention) == 1 {
			return nil, fmt.Errorf("error, this: %s is not an extention", s)
		}

		// add to map
		m[extention] = true

	}

	return m, nil
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
