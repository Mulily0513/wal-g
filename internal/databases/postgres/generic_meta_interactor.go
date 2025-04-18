package postgres

import (
	"github.com/pkg/errors"
	"github.com/wal-g/wal-g/internal"
	"github.com/wal-g/wal-g/pkg/storages/storage"
)

type GenericMetaInteractor struct {
	GenericMetaFetcher
	GenericMetaSetter
}

func NewGenericMetaInteractor() GenericMetaInteractor {
	return GenericMetaInteractor{
		GenericMetaFetcher: NewGenericMetaFetcher(),
		GenericMetaSetter:  NewGenericMetaSetter(),
	}
}

type GenericMetaFetcher struct{}

func NewGenericMetaFetcher() GenericMetaFetcher {
	return GenericMetaFetcher{}
}

func (mf GenericMetaFetcher) Fetch(backupName string, backupFolder storage.Folder) (internal.GenericMetadata, error) {
	backup, err := NewBackup(backupFolder, backupName)
	if err != nil {
		return internal.GenericMetadata{}, err
	}
	meta, err := backup.FetchMeta()
	if err != nil {
		return internal.GenericMetadata{}, err
	}

	return internal.GenericMetadata{
		BackupName:       backupName,
		UncompressedSize: meta.UncompressedSize,
		CompressedSize:   meta.CompressedSize,
		Hostname:         meta.Hostname,
		StartTime:        meta.StartTime,
		FinishTime:       meta.FinishTime,
		IsPermanent:      meta.IsPermanent,
		IncrementDetails: NewIncrementDetailsFetcher(backup),
		UserData:         meta.UserData,
	}, nil
}

// TODO rewrite multistorage pg operations with this method and internal.GetPermanentBackups instead of postgres.GetPermanentBackupsAndWals
func (mf GenericMetaFetcher) FetchFromStorage(
	backupName string, backupFolder storage.Folder, storage string,
) (internal.GenericMetadata, error) {
	return mf.Fetch(backupName, backupFolder)
}

type GenericMetaSetter struct{}

func NewGenericMetaSetter() GenericMetaSetter {
	return GenericMetaSetter{}
}

func (ms GenericMetaSetter) SetUserData(backupName string, backupFolder storage.Folder, userData interface{}) error {
	modifier := func(dto ExtendedMetadataDto) ExtendedMetadataDto {
		dto.UserData = userData
		return dto
	}
	return modifyBackupMetadata(backupName, backupFolder, modifier)
}

func (ms GenericMetaSetter) SetIsPermanent(backupName string, backupFolder storage.Folder, isPermanent bool) error {
	modifier := func(dto ExtendedMetadataDto) ExtendedMetadataDto {
		dto.IsPermanent = isPermanent
		return dto
	}
	return modifyBackupMetadata(backupName, backupFolder, modifier)
}

func modifyBackupMetadata(backupName string, backupFolder storage.Folder, modifier func(ExtendedMetadataDto) ExtendedMetadataDto) error {
	backup, err := internal.NewBackup(backupFolder, backupName)
	if err != nil {
		return errors.Wrap(err, "failed to modify metadata")
	}
	var meta ExtendedMetadataDto
	err = backup.FetchMetadata(&meta)
	if err != nil {
		return errors.Wrap(err, "failed to fetch the existing backup metadata for modifying")
	}
	meta = modifier(meta)
	err = backup.UploadMetadata(meta)
	if err != nil {
		return errors.Wrap(err, "failed to upload the modified metadata to the storage")
	}
	return nil
}

type IncrementDetailsFetcher struct {
	backup Backup
}

func NewIncrementDetailsFetcher(backup Backup) *IncrementDetailsFetcher {
	return &IncrementDetailsFetcher{backup: backup}
}

func (idf *IncrementDetailsFetcher) Fetch() (bool, internal.IncrementDetails, error) {
	sentinel, err := idf.backup.GetSentinel()
	if err != nil || !sentinel.IsIncremental() {
		return false, internal.IncrementDetails{}, err
	}

	return true, internal.IncrementDetails{
		IncrementFrom:     *sentinel.IncrementFrom,
		IncrementFullName: *sentinel.IncrementFullName,
		IncrementCount:    *sentinel.IncrementCount,
	}, nil
}
