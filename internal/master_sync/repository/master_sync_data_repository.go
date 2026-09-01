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

func (repo *MasterSyncRepository) Create(
	ctx context.Context,
	req domain.DataMasterSync,
) error {
	const op = "master-sync-repo-Create"

	repo.Logger.Debug(
		"create sync record",
		"op", op,
		"id", req.ID,
		"sync_type", req.SyncType,
		"sync_from", req.SyncFrom,
		"status", req.Status,
	)

	const query = `
		INSERT INTO data_master_sync (
			id,
			sync_from,
			sync_type,
			started_at,
			status,
			total_data,
			total_success,
			total_failed
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := repo.DB.ExecContext(
		ctx,
		query,
		req.ID,
		req.SyncFrom,
		req.SyncType,
		req.StartedAt,
		req.Status,
		req.TotalData,
		req.TotalSuccess,
		req.TotalFailed,
	)
	if err != nil {
		repo.Logger.Error(
			"failed to create sync record",
			"op", op,
			"id", req.ID,
			"error", err,
		)
		return err
	}

	repo.Logger.Debug(
		"successfully created sync record",
		"op", op,
		"id", req.ID,
	)

	return nil
}

func (repo *MasterSyncRepository) Finish(
	ctx context.Context,
	req domain.DataMasterSync,
) error {
	const op = "master-sync-repo-Finish"

	repo.Logger.Debug(
		"update sync record",
		"op", op,
		"id", req.ID,
		"status", req.Status,
	)

	const query = `
		UPDATE data_master_sync
		SET
			finished_at = $1,
			status = $2,
			total_data = $3,
			total_success = $4,
			total_failed = $5
		WHERE id = $6
	`

	_, err := repo.DB.ExecContext(
		ctx,
		query,
		req.FinishedAt,
		req.Status,
		req.TotalData,
		req.TotalSuccess,
		req.TotalFailed,
		req.ID,
	)
	if err != nil {
		repo.Logger.Error(
			"failed to update sync record",
			"op", op,
			"id", req.ID,
			"error", err,
		)
		return err
	}

	repo.Logger.Debug(
		"successfully updated sync record",
		"op", op,
		"id", req.ID,
	)

	return nil
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
