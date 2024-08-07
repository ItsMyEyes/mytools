package usecase

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/api/drive/v3"
)

type BackupUsecase interface {
	BackupSixHours(ctx context.Context, req RequestBackupSixHours) error
}

type backupUseCase struct {
	driveServices *drive.Service
}

func NewBackupUsecase(driveServices *drive.Service) BackupUsecase {
	return backupUseCase{
		driveServices: driveServices,
	}
}

type RequestBackupSixHours struct {
	FilePath       string
	RemoveBackup   bool
	HowOldDuration time.Duration
	Parents        []string
}

func (u backupUseCase) BackupSixHours(ctx context.Context, req RequestBackupSixHours) error {
	f, err := os.Open(req.FilePath)
	if err != nil {
		err = fmt.Errorf("[BackupUseCase] Path %s Failed to open file: %w", req.FilePath, err)
		log.Error().Err(err).Msg("BackupUseCase")
		return err
	}
	defer f.Close()

	fileModelsGDrive := &drive.File{
		Name:    fmt.Sprintf("%s:%s", time.Now().Format("2006-01-02@15:04"), f.Name()),
		Parents: req.Parents,
	}

	_, err = u.driveServices.Files.Create(fileModelsGDrive).Media(f).Do()
	if err != nil {
		err = fmt.Errorf("[BackupUseCase] Failed to upload file to GDrive: %w", err)
		log.Error().Err(err).Msg("BackupUseCase")
		return err
	}

	if req.RemoveBackup {
		name := fmt.Sprintf("%s:%s", time.Now().Add(req.HowOldDuration).Format("2006-01-02@15:04"), f.Name())
		log.Info().Msgf("Looking for old backup file: %s", name)
		oldBackupFile, err := u.FindGDriveByName(ctx)
		if err != nil {
			err = fmt.Errorf("[BackupUseCase] Failed to find old backup file: %w", err)
			log.Error().Err(err).Msg("BackupUseCase")
			return err
		}
		for _, reply := range oldBackupFile {
			if reply.Name == name {
				log.Info().Msgf("Deleting old backup file %s", reply.Name)
				err := u.driveServices.Files.Delete(reply.Id).Do()
				if err != nil {
					err = fmt.Errorf("[BackupUseCase] Failed to delete old backup file: %w", err)
					log.Error().Err(err).Msg("BackupUseCase")
					return err
				}
			}
		}
	}

	return nil
}

func (u backupUseCase) FindGDriveByName(ctx context.Context) ([]*drive.File, error) {
	files, err := u.driveServices.Files.List().Do()
	if err != nil {
		err = fmt.Errorf("[FindGDriveByName] Failed to find GDrive file: %w", err)
		return nil, err
	}

	return files.Files, nil
}
