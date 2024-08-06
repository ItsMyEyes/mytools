package cli

import (
	"context"
	"fmt"
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
			fmt.Println("Args: ", args)
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
			initUseCase.BackupSixHours(ctx, usecase.RequestBackupSixHours{
				FilePath:       args[0],
				RemoveBackup:   cmd.Flag("remove-backup").Value.String() == "true",
				HowOldDuration: time.Duration(-durationInt) * time.Second,
				Parents:        []string{cmd.Flag("parents").Value.String()},
			})
		},
	}
	rootCmd.Flags().String("path", "", "file path to backup")
	rootCmd.Flags().Int("duration", 0, "how old duration in second for remove backup in drive")
	rootCmd.Flags().String("parents", "", "parents for backup")
	rootCmd.Flags().Bool("remove-backup", false, "RemoveBackup")

	return command
}
