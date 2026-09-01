package repository

import (
	"context"
	"database/sql"
	"log/slog"

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
		WHERE sync_id = ?
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
