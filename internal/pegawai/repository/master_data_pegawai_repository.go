package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"master-data-service/internal/pegawai/domain"
)

type MasterDataPegawaiRepository struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewMasterDataPegawaiRepository(
	db *sql.DB,
	logger *slog.Logger,
) *MasterDataPegawaiRepository {
	return &MasterDataPegawaiRepository{
		DB:     db,
		Logger: logger,
	}
}

func (repo *MasterDataPegawaiRepository) BulkInsert(
	ctx context.Context,
	pegawais []domain.DataMasterPegawai,
	syncID string,
	syncFrom string,
) (int, error) {
	const op = "pegawai-repo-BulkInsert"

	repo.Logger.Debug(
		"bulk insert pegawai",
		"op", op,
		"sync_id", syncID,
		"sync_from", syncFrom,
		"count", len(pegawais),
	)

	const query = `
		INSERT INTO data_master_pegawai (
			opd_kode,

			pegawai_nip,
			pegawai_nama,
			pegawai_jenis_kelamin,

			pegawai_lahir_tempat,
			pegawai_lahir_tanggal,
			pegawai_pensiun_tmt,

			pegawai_aktif_id,
			pegawai_status,
			pegawai_jenis,

			pegawai_golongan_kode,
			pegawai_golongan_nama,
			pegawai_golongan_pangkat,
			pegawai_golongan_tmt,

			pegawai_jenis_jabatan_kode,
			pegawai_jenis_jabatan_nama,

			pegawai_jabatan_kode,
			pegawai_jabatan_terakhir,

			pegawai_eselon_kode,
			pegawai_eselon_nama,

			opd_level,
			opd_nama,

			sync_id,
			sync_from,
			sync_at
		)
		VALUES (
			$1,

			$2,
			$3,
			$4,

			$5,
			$6,
			$7,

			$8,
			$9,
			$10,

			$11,
			$12,
			$13,
			$14,

			$15,
			$16,

			$17,
			$18,

			$19,
			$20,

			$21,
			$22,

			$23,
			$24,
			$25
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

	for _, pegawai := range pegawais {
		result, err := stmt.ExecContext(
			ctx,

			pegawai.OpdKode,

			pegawai.PegawaiNIP,
			pegawai.PegawaiNama,
			pegawai.PegawaiJenisKelamin,

			pegawai.PegawaiLahirTempat,
			pegawai.PegawaiLahirTanggal,
			pegawai.PegawaiPensiunTMT,

			pegawai.PegawaiAktifID,
			pegawai.PegawaiStatus,
			pegawai.PegawaiJenis,

			pegawai.PegawaiGolonganKode,
			pegawai.PegawaiGolonganNama,
			pegawai.PegawaiGolonganPangkat,
			pegawai.PegawaiGolonganTMT,

			pegawai.PegawaiJenisJabatanKode,
			pegawai.PegawaiJenisJabatanNama,

			pegawai.PegawaiJabatanKode,
			pegawai.PegawaiJabatanTerakhir,

			pegawai.PegawaiEselonKode,
			pegawai.PegawaiEselonNama,

			pegawai.OpdLevel,
			pegawai.OpdNama,

			syncID,
			syncFrom,
			now,
		)
		if err != nil {
			repo.Logger.Error(
				"failed to insert pegawai",
				"op", op,
				"sync_id", syncID,
				"opd_kode", pegawai.OpdKode,
				"pegawai_nip", pegawai.PegawaiNIP,
				"error", err,
			)

			return inserted, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return inserted, err
		}

		if rowsAffected == 1 {
			inserted++
			continue
		}

		repo.Logger.Debug(
			"pegawai already exists, skip",
			"op", op,
			"sync_id", syncID,
			"pegawai_nip", pegawai.PegawaiNIP,
		)

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
		"successfully bulk inserted pegawai",
		"op", op,
		"sync_id", syncID,
		"inserted", inserted,
	)

	return inserted, nil
}

func (repo *MasterDataPegawaiRepository) FindAllPegawaiBySnapshotIdAndKodeOpd(
	ctx context.Context,
	snapshotID string,
	kodeOpd string,
) ([]domain.DataMasterPegawai, error) {
	const op = "pegawai-repo-FindAllPegawaiBySnapshotIdAndKodeOpd"

	repo.Logger.Debug(
		"find all pegawai by snapshot and opd kode",
		"op", op,
		"snapshot_id", snapshotID,
		"opd_kode", kodeOpd,
	)

	const query = `
		SELECT
			id,
			opd_kode,

			pegawai_nip,
			pegawai_nama,
			pegawai_jenis_kelamin,

			pegawai_lahir_tempat,
			pegawai_lahir_tanggal,
			pegawai_pensiun_tmt,

			pegawai_aktif_id,
			pegawai_status,
			pegawai_jenis,

			pegawai_golongan_kode,
			pegawai_golongan_nama,
			pegawai_golongan_pangkat,
			pegawai_golongan_tmt,

			pegawai_jenis_jabatan_kode,
			pegawai_jenis_jabatan_nama,

			pegawai_jabatan_kode,
			pegawai_jabatan_terakhir,

			pegawai_eselon_kode,
			pegawai_eselon_nama,

			opd_level,
			opd_nama,

			sync_id,
			sync_from,
			sync_at,
			created_at,
			updated_at
		FROM data_master_pegawai
		WHERE sync_id = $1
		  AND opd_kode = $2
		ORDER BY pegawai_nama
	`

	rows, err := repo.DB.QueryContext(
		ctx,
		query,
		snapshotID,
		kodeOpd,
	)
	if err != nil {
		repo.Logger.Error(
			"failed to query pegawai",
			"op", op,
			"snapshot_id", snapshotID,
			"opd_kode", kodeOpd,
			"error", err,
		)
		return nil, err
	}
	defer rows.Close()

	pegawais := make([]domain.DataMasterPegawai, 0)

	for rows.Next() {
		var pegawai domain.DataMasterPegawai

		err := rows.Scan(
			&pegawai.ID,
			&pegawai.OpdKode,

			&pegawai.PegawaiNIP,
			&pegawai.PegawaiNama,
			&pegawai.PegawaiJenisKelamin,

			&pegawai.PegawaiLahirTempat,
			&pegawai.PegawaiLahirTanggal,
			&pegawai.PegawaiPensiunTMT,

			&pegawai.PegawaiAktifID,
			&pegawai.PegawaiStatus,
			&pegawai.PegawaiJenis,

			&pegawai.PegawaiGolonganKode,
			&pegawai.PegawaiGolonganNama,
			&pegawai.PegawaiGolonganPangkat,
			&pegawai.PegawaiGolonganTMT,

			&pegawai.PegawaiJenisJabatanKode,
			&pegawai.PegawaiJenisJabatanNama,

			&pegawai.PegawaiJabatanKode,
			&pegawai.PegawaiJabatanTerakhir,

			&pegawai.PegawaiEselonKode,
			&pegawai.PegawaiEselonNama,

			&pegawai.OpdLevel,
			&pegawai.OpdNama,

			&pegawai.SyncID,
			&pegawai.SyncFrom,
			&pegawai.SyncAt,
			&pegawai.CreatedAt,
			&pegawai.UpdatedAt,
		)
		if err != nil {
			repo.Logger.Error(
				"failed to scan pegawai",
				"op", op,
				"snapshot_id", snapshotID,
				"opd_kode", kodeOpd,
				"error", err,
			)
			return nil, err
		}

		pegawais = append(pegawais, pegawai)
	}

	if err := rows.Err(); err != nil {
		repo.Logger.Error(
			"error while iterating pegawai rows",
			"op", op,
			"snapshot_id", snapshotID,
			"opd_kode", kodeOpd,
			"error", err,
		)
		return nil, err
	}

	repo.Logger.Debug(
		"successfully found pegawai",
		"op", op,
		"snapshot_id", snapshotID,
		"opd_kode", kodeOpd,
		"count", len(pegawais),
	)

	return pegawais, nil
}
