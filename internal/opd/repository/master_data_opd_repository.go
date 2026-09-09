package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"master-data-service/internal/opd/domain"
)

type MasterDataOpdRepository struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewMasterDataOpdRepository(
	db *sql.DB,
	logger *slog.Logger,
) *MasterDataOpdRepository {
	return &MasterDataOpdRepository{
		DB:     db,
		Logger: logger,
	}
}

func (repo *MasterDataOpdRepository) BulkInsert(
	ctx context.Context,
	opds []domain.DataMasterOpd,
	syncID string,
	syncFrom string,
) (int, error) {
	const op = "opd-repo-BulkInsert"

	repo.Logger.Debug(
		"bulk insert opd",
		"op", op,
		"sync_id", syncID,
		"sync_from", syncFrom,
		"count", len(opds),
	)

	const query = `
		INSERT INTO data_master_opd (
			opd_kode,
			opd_level,
			opd_nama,
			opd_jabatan,
			eselon_kode,
			eselon_nama,
			sync_id,
			sync_from,
			sync_at
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9
		)
	`

	now := time.Now()

	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		repo.Logger.Error(
			"failed to begin transaction",
			"op", op,
			"error", err,
		)
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		repo.Logger.Error(
			"failed to prepare insert statement",
			"op", op,
			"error", err,
		)
		return 0, err
	}
	defer stmt.Close()

	inserted := 0

	for _, opd := range opds {
		_, err := stmt.ExecContext(
			ctx,
			opd.OpdKode,
			opd.OpdLevel,
			opd.OpdNama,
			opd.OpdJabatan,
			opd.EselonKode,
			opd.EselonNama,
			syncID,
			syncFrom,
			now,
		)
		if err != nil {
			repo.Logger.Error(
				"failed to insert opd",
				"op", op,
				"sync_id", syncID,
				"opd_kode", opd.OpdKode,
				"error", err,
			)

			return inserted, err
		}

		inserted++
	}

	if err := tx.Commit(); err != nil {
		repo.Logger.Error(
			"failed to commit transaction",
			"op", op,
			"sync_id", syncID,
			"error", err,
		)
		return inserted, err
	}

	repo.Logger.Debug(
		"successfully bulk inserted opd",
		"op", op,
		"sync_id", syncID,
		"inserted", inserted,
	)

	return inserted, nil
}

func (repo *MasterDataOpdRepository) FindAllOpdBySnapshotId(
	ctx context.Context,
	snapshotId string,
) ([]domain.DataMasterOpd, error) {
	const op = "opd-repo-FindAllOpdBySnapshotId"

	repo.Logger.Debug(
		"find all opd by snapshot",
		"op", op,
		"snapshot_id", snapshotId,
	)

	const query = `
		SELECT
			id,
			opd_kode,
			opd_level,
			opd_nama,
			opd_jabatan,
			eselon_kode,
			eselon_nama,
			sync_id,
			sync_from,
			sync_at,
			created_at,
			updated_at
		FROM data_master_opd
		WHERE sync_id = $1
		AND opd_level = 1
		ORDER BY opd_kode
	`

	rows, err := repo.DB.QueryContext(ctx, query, snapshotId)
	if err != nil {
		repo.Logger.Error(
			"failed to query opd",
			"op", op,
			"snapshot_id", snapshotId,
			"error", err,
		)
		return nil, err
	}
	defer rows.Close()

	opds := make([]domain.DataMasterOpd, 0)

	for rows.Next() {
		var opd domain.DataMasterOpd

		err := rows.Scan(
			&opd.ID,
			&opd.OpdKode,
			&opd.OpdLevel,
			&opd.OpdNama,
			&opd.OpdJabatan,
			&opd.EselonKode,
			&opd.EselonNama,
			&opd.SyncID,
			&opd.SyncFrom,
			&opd.SyncAt,
			&opd.CreatedAt,
			&opd.UpdatedAt,
		)
		if err != nil {
			repo.Logger.Error(
				"failed to scan opd",
				"op", op,
				"snapshot_id", snapshotId,
				"error", err,
			)
			return nil, err
		}

		opds = append(opds, opd)
	}

	if err := rows.Err(); err != nil {
		repo.Logger.Error(
			"error while iterating opd rows",
			"op", op,
			"snapshot_id", snapshotId,
			"error", err,
		)
		return nil, err
	}

	repo.Logger.Debug(
		"successfully found opd",
		"op", op,
		"snapshot_id", snapshotId,
		"count", len(opds),
	)

	return opds, nil
}

func (repo *MasterDataOpdRepository) FindOpdByKode(
	ctx context.Context,
	snapshotId string,
	kodeOpd string,
) (domain.DataMasterOpd, error) {
	const op = "opd-repo-FindOpdByKode"

	repo.Logger.Debug(
		"find opd by kode",
		"op", op,
		"snapshot_id", snapshotId,
		"opd_kode", kodeOpd,
	)

	const query = `
		SELECT
			id,
			opd_kode,
			opd_level,
			opd_nama,
			opd_jabatan,
			eselon_kode,
			eselon_nama,
			sync_id,
			sync_from,
			sync_at,
			created_at,
			updated_at
		FROM data_master_opd
		WHERE sync_id = $1
		  AND opd_kode = $2
		LIMIT 1
	`

	var opd domain.DataMasterOpd

	err := repo.DB.QueryRowContext(
		ctx,
		query,
		snapshotId,
		kodeOpd,
	).Scan(
		&opd.ID,
		&opd.OpdKode,
		&opd.OpdLevel,
		&opd.OpdNama,
		&opd.OpdJabatan,
		&opd.EselonKode,
		&opd.EselonNama,
		&opd.SyncID,
		&opd.SyncFrom,
		&opd.SyncAt,
		&opd.CreatedAt,
		&opd.UpdatedAt,
	)

	if err != nil {
		repo.Logger.Error(
			"failed to find opd",
			"op", op,
			"snapshot_id", snapshotId,
			"opd_kode", kodeOpd,
			"error", err,
		)
		return opd, err
	}

	repo.Logger.Debug(
		"successfully found opd",
		"op", op,
		"snapshot_id", snapshotId,
		"opd_kode", kodeOpd,
		"opd_id", opd.ID,
	)

	return opd, nil
}
