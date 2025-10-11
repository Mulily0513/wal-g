package gp

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wal-g/tracelog"
	"github.com/wal-g/wal-g/internal"
	conf "github.com/wal-g/wal-g/internal/config"
	"github.com/wal-g/wal-g/internal/databases/postgres"
	"github.com/wal-g/wal-g/utility"
)

const (
	segBackupMergeShortDescription = "Merges incremental backups for a single segment"
)

var (
	segContentID int

	// segBackupMergeCmd is a subcommand to merge incremental backups of a single segment.
	// It is called remotely by a backup-merge command from the master host
	segBackupMergeCmd = &cobra.Command{
		Use:   "seg-backup-merge target_backup_name --content-id=[content_id]",
		Short: segBackupMergeShortDescription,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			targetBackupName := args[0]
			storage, err := internal.ConfigureStorage()
			tracelog.ErrorLogger.FatalOnError(err)
			folder := storage.RootFolder()

			tarBallComposerType := chooseTarBallComposer()

			// Use the basebackups subdirectory for backup operations
			folder = folder.GetSubFolder(utility.SegmentsPath + fmt.Sprintf("/seg%d/", segContentID))

			// Create PostgreSQL backup merge handler for this segment
			mergeHandler, err := postgres.NewBackupMergeHandler(targetBackupName, folder, tarBallComposerType)
			tracelog.ErrorLogger.FatalOnError(err)

			mergeHandler.HandleBackupMerge()
		},
	}
)

func chooseTarBallComposer() postgres.TarBallComposerType {
	tarBallComposerType := postgres.RegularComposer

	useRatingComposer := viper.GetBool(conf.UseRatingComposerSetting)
	if useRatingComposer {
		tarBallComposerType = postgres.RatingComposer
	}

	useDatabaseComposer := viper.GetBool(conf.UseDatabaseComposerSetting)
	if useDatabaseComposer {
		tarBallComposerType = postgres.DatabaseComposer
	}

	useCopyComposer := viper.GetBool(conf.UseCopyComposerSetting)
	if useCopyComposer {
		fullBackup = true
		tarBallComposerType = postgres.CopyComposer
	}

	return tarBallComposerType
}

func init() {
	// Since this is a utility command, it should not be exposed to the end user.
	segBackupMergeCmd.Hidden = true
	segBackupMergeCmd.PersistentFlags().IntVar(&segContentID, "content-id", 0, "segment content ID")
	_ = segBackupMergeCmd.MarkFlagRequired("content-id")
	cmd.AddCommand(segBackupMergeCmd)
}
