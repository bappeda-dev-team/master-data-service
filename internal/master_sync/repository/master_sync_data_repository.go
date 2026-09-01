package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"master-data-service/internal/master_sync/domain"
)

type MasterSyncRepository struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewMasterSyncRepository(
	db *sql.DB,
	logger *slog.Logger,
) *MasterSyncRepository {
	return &MasterSyncRepository{
		DB:     db,
		Logger: logger,
	}
}

func (repo *MasterSyncRepository) FindActiveSnapshot(
	ctx context.Context,
	req domain.DataMasterSync,
) (string, error) {
	const op = "master-sync-repo-FindActiveSnapshot"

	repo.Logger.Debug(
		"find active snapshot",
		"op", op,
		"sync_type", req.SyncType,
		"sync_from", req.SyncFrom,
		"status", req.Status,
	)
	query := `SELECT id
		FROM data_master_sync
		WHERE sync_type = $1
		AND sync_from = $2
		AND status = $3`
	var id string
	err := repo.DB.QueryRowContext(
		ctx,
		query,
		req.SyncType,
		req.SyncFrom,
		req.Status,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
