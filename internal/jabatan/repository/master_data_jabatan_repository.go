package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"master-data-service/internal/jabatan/domain"
)

type MasterDataJabatanRepository struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewMasterDataJabatanRepository(
	db *sql.DB,
	logger *slog.Logger,
) *MasterDataJabatanRepository {
	return &MasterDataJabatanRepository{
		DB:     db,
		Logger: logger,
	}
}

func (repo *MasterDataJabatanRepository) BulkInsert(
	ctx context.Context,
	jabatans []domain.DataMasterJabatan,
	syncID string,
	syncFrom string,
) (int, error) {
	const op = "jabatan-repo-BulkInsert"

	repo.Logger.Debug(
		"bulk insert jabatan",
		"op", op,
		"sync_id", syncID,
		"sync_from", syncFrom,
		"count", len(jabatans),
	)

	const query = `
		INSERT INTO data_master_jabatan (
			opd_kode,
			jenis_jabatan_kode,
			jenis_jabatan_nama,
			jabatan_kode,
			jabatan_nama,
			jenjang_jabatan_kode,
			jenjang_jabatan_nama,
			sync_id,
			sync_from,
			sync_at
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10
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

	for _, jabatan := range jabatans {
		_, err := stmt.ExecContext(
			ctx,
			jabatan.OpdKode,
			jabatan.JenisJabatanKode,
			jabatan.JenisJabatanNama,
			jabatan.JabatanKode,
			jabatan.JabatanNama,
			jabatan.JenjangJabatanKode,
			jabatan.JenjangJabatanNama,
			syncID,
			syncFrom,
			now,
		)
		if err != nil {
			repo.Logger.Error(
				"failed to insert jabatan",
				"op", op,
				"opd_kode", jabatan.OpdKode,
				"jenis_jabatan_kode", jabatan.JenisJabatanKode,
				"jabatan_kode", jabatan.JabatanKode,
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
			"error", err,
		)
		return inserted, err
	}

	repo.Logger.Debug(
		"successfully bulk inserted jabatan",
		"op", op,
		"sync_id", syncID,
		"inserted", inserted,
	)

	return inserted, nil
}

func (repo *MasterDataJabatanRepository) FindAllJabatanBySnapshotId(
	ctx context.Context,
	snapshotID string,
) ([]domain.DataMasterJabatan, error) {
	const op = "jabatan-repo-FindAllJabatanBySnapshotId"

	repo.Logger.Debug(
		"find all jabatan by snapshot",
		"op", op,
		"snapshot_id", snapshotID,
	)

	const query = `
		SELECT
			id,
			opd_kode,
			jenis_jabatan_kode,
			jenis_jabatan_nama,
			jabatan_kode,
			jabatan_nama,
			jenjang_jabatan_kode,
			jenjang_jabatan_nama,
			sync_id,
			sync_from,
			sync_at,
			created_at,
			updated_at
		FROM data_master_jabatan
		WHERE sync_id = $1
		ORDER BY opd_kode, jenis_jabatan_kode, jabatan_kode
	`

	rows, err := repo.DB.QueryContext(ctx, query, snapshotID)
	if err != nil {
		repo.Logger.Error(
			"failed to query jabatan",
			"op", op,
			"snapshot_id", snapshotID,
			"error", err,
		)
		return nil, err
	}
	defer rows.Close()

	jabatans := make([]domain.DataMasterJabatan, 0)

	for rows.Next() {
		var jabatan domain.DataMasterJabatan

		err := rows.Scan(
			&jabatan.ID,
			&jabatan.OpdKode,
			&jabatan.JenisJabatanKode,
			&jabatan.JenisJabatanNama,
			&jabatan.JabatanKode,
			&jabatan.JabatanNama,
			&jabatan.JenjangJabatanKode,
			&jabatan.JenjangJabatanNama,
			&jabatan.SyncID,
			&jabatan.SyncFrom,
			&jabatan.SyncAt,
			&jabatan.CreatedAt,
			&jabatan.UpdatedAt,
		)
		if err != nil {
			repo.Logger.Error(
				"failed to scan jabatan",
				"op", op,
				"snapshot_id", snapshotID,
				"error", err,
			)
			return nil, err
		}

		jabatans = append(jabatans, jabatan)
	}

	if err := rows.Err(); err != nil {
		repo.Logger.Error(
			"error while iterating jabatan rows",
			"op", op,
			"snapshot_id", snapshotID,
			"error", err,
		)
		return nil, err
	}

	repo.Logger.Debug(
		"successfully found jabatan",
		"op", op,
		"snapshot_id", snapshotID,
		"count", len(jabatans),
	)

	return jabatans, nil
}

func (repo *MasterDataJabatanRepository) FindAllJabatanBySnapshotIdAndKodeOpd(
	ctx context.Context,
	snapshotID string,
	kodeOpd string,
) ([]domain.DataMasterJabatan, error) {
	const op = "jabatan-repo-FindAllJabatanBySnapshotIdAndKodeOpd"

	repo.Logger.Debug(
		"find all jabatan by snapshot and opd kode",
		"op", op,
		"snapshot_id", snapshotID,
		"opd_kode", kodeOpd,
	)

	const query = `
		SELECT
			id,
			opd_kode,
			jenis_jabatan_kode,
			jenis_jabatan_nama,
			jabatan_kode,
			jabatan_nama,
			jenjang_jabatan_kode,
			jenjang_jabatan_nama,
			sync_id,
			sync_from,
			sync_at,
			created_at,
			updated_at
		FROM data_master_jabatan
		WHERE sync_id = $1
		  AND opd_kode = $2
		ORDER BY jabatan_nama
	`

	rows, err := repo.DB.QueryContext(
		ctx,
		query,
		snapshotID,
		kodeOpd,
	)
	if err != nil {
		repo.Logger.Error(
			"failed to query jabatan",
			"op", op,
			"snapshot_id", snapshotID,
			"opd_kode", kodeOpd,
			"error", err,
		)
		return nil, err
	}
	defer rows.Close()

	jabatans := make([]domain.DataMasterJabatan, 0)

	for rows.Next() {
		var jabatan domain.DataMasterJabatan

		err := rows.Scan(
			&jabatan.ID,
			&jabatan.OpdKode,
			&jabatan.JenisJabatanKode,
			&jabatan.JenisJabatanNama,
			&jabatan.JabatanKode,
			&jabatan.JabatanNama,
			&jabatan.JenjangJabatanKode,
			&jabatan.JenjangJabatanNama,
			&jabatan.SyncID,
			&jabatan.SyncFrom,
			&jabatan.SyncAt,
			&jabatan.CreatedAt,
			&jabatan.UpdatedAt,
		)
		if err != nil {
			repo.Logger.Error(
				"failed to scan jabatan",
				"op", op,
				"snapshot_id", snapshotID,
				"opd_kode", kodeOpd,
				"error", err,
			)
			return nil, err
		}

		jabatans = append(jabatans, jabatan)
	}

	if err := rows.Err(); err != nil {
		repo.Logger.Error(
			"error while iterating jabatan rows",
			"op", op,
			"snapshot_id", snapshotID,
			"opd_kode", kodeOpd,
			"error", err,
		)
		return nil, err
	}

	repo.Logger.Debug(
		"successfully found jabatan",
		"op", op,
		"snapshot_id", snapshotID,
		"opd_kode", kodeOpd,
		"count", len(jabatans),
	)

	return jabatans, nil
}
