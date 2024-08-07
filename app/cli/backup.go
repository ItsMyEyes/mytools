package cli

import (
	"context"
	"strconv"
	"time"

	drive_google "github.com/ItsMyEyes/my-tools/internal/lib/google/drive"
	"github.com/ItsMyEyes/my-tools/internal/usecase"
	"github.com/spf13/cobra"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func NewBackupCLI() *cobra.Command {
	command := &cobra.Command{
		Use:        "backup",
		Aliases:    []string{"b"},
		Short:      "run backup for database",
		ArgAliases: []string{"a"},
		Run: func(cmd *cobra.Command, args []string) {
			// get flag value
			authSrv := drive_google.AuthService{}
			ctx := context.TODO()
			client, err := authSrv.GetGoogleClient(ctx, "./credentials/credentials.json")
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get google client")
			}

			if len(args) < 1 {
				log.Fatal().Msg("File path is required")
			}

			// Create Google Drive Service
			srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get google drive service")
			}
			durationInt, _ := strconv.Atoi(cmd.Flag("duration").Value.String())
			initUseCase := usecase.NewBackupUsecase(srv)
			initUseCase.BackupNow(ctx, usecase.RequestBackupNow{
				FilePath:       args[0],
				RemoveBackup:   cmd.Flag("remove-backup").Value.String() == "true",
				HowOldDuration: time.Duration(-durationInt) * time.Hour,
				Parents:        []string{cmd.Flag("parents").Value.String()},
			})
		},
	}
	command.Flags().String("path", "", "file path to backup")
	command.Flags().Int("duration", 1, "how old duration in Hour for remove backup in drive")
	command.Flags().String("parents", "", "parents for backup")
	command.Flags().Bool("remove-backup", false, "RemoveBackup")

	return command
}
