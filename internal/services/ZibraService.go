package services

import (
	"fmt"
	"os"

	"github.com/Talos-hub/ZibraGo/internal/apperrors"
	"github.com/Talos-hub/ZibraGo/internal/ports"
)

type ZibraService struct {
	walker  ports.Walker
	arhiver ports.Archiver
	api     ports.ApiCloud
	logger  ports.Logger
}

// NewZibra returns pointer to ZibraService
func NewZibra(walker ports.Walker, arhiver ports.Archiver, api ports.ApiCloud, logger ports.Logger) (*ZibraService, error) {
	// check interfaces
	if walker == nil || arhiver == nil || api == nil {
		return nil, apperrors.NewAppError("expected non-nil parameters", "NewZibra", apperrors.E_FATAL, nil)
	}
	return &ZibraService{
		walker:  walker,
		arhiver: arhiver,
		api:     api,
		logger:  logger,
	}, nil
}

// Run starts work service.
func (z *ZibraService) Run() error {

	// check connect to api
	err := z.api.Check()
	if err != nil {
		z.logger.Error("Cloud API unavailable", "error", err)
		return fmt.Errorf("cloud API unavailable: %w", err)
	}

	// scan dir
	pathes, err := z.walker.Walk()
	if err != nil {
		z.logger.Error("Error scanning directory", "error", err)
		return fmt.Errorf("scan error: %w", err)
	}

	if len(pathes) == 0 {
		z.logger.Warn("Directory is empty", "Dir", pathes)
		return apperrors.NewAppError("dirictory is empty", "Run", apperrors.E_EMPTY_DIR, nil)
	}

	z.logger.Info("Files found", "count", len(pathes))

	// start archiving files
	filepath, err := z.arhiver.Start(pathes...)
	if err != nil {
		z.logger.Error("Error start zip archiver", "error", err)
		return fmt.Errorf("error run zibra service: %w", err)
	}
	// open a zip file
	zipFile, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			z.logger.Error("Fatal error, zip file is not exist", "error", err)
			return apperrors.NewAppError("error, zip file is not exit", "Run", apperrors.E_FATAL, err)
		}
		z.logger.Error("Error open a zip file", "error", err)
		return apperrors.NewAppError("error open zip file", "Run", apperrors.E_OPEN, err)
	}
	defer zipFile.Close()

	err = z.api.UpLoadFile(zipFile)
	if err != nil {
		z.logger.Error("Upload failed", "error", err)
		return fmt.Errorf("upload error: %w", err)
	}

	return nil
}
