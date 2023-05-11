package pg

import (
	"github.com/spf13/cobra"
	"github.com/wal-g/tracelog"
	"github.com/wal-g/wal-g/internal"
	"github.com/wal-g/wal-g/internal/databases/postgres"
)

var targetIncrementalBackupName string

const (
	backupMergeShortDescription            = "Create a single backup from delta backups and put it in storage"
	targetIncrementalBackupNameDescription = "Name of the target delta backup relative to which the base backup should be generated"
)

var backupMergeCmd = &cobra.Command{
	Use:   "backup-merge backup_name",
	Short: backupMergeShortDescription,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetBackupName := args[0]
		storage, err := internal.ConfigureStorage()
		tracelog.ErrorLogger.FatalOnError(err)
		folder := storage.RootFolder()

		composer := chooseTarBallComposer()

		mergeHandler, err := postgres.NewBackupMergeHandler(targetBackupName, folder, composer)
		tracelog.ErrorLogger.FatalOnError(err)

		mergeHandler.HandleBackupMerge()
		tracelog.InfoLogger.Println("DONE")
	},
}

func init() {
	backupMergeCmd.Flags().StringVar(&targetIncrementalBackupName, "target-backup-name", "",
		targetIncrementalBackupNameDescription)

	Cmd.AddCommand(backupMergeCmd)
}
