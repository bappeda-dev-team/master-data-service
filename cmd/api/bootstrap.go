package main

import (
	"database/sql"
	"log/slog"

	"master-data-service/internal/api"

	masterSyncRepository "master-data-service/internal/master_sync/repository"
	masterSyncService "master-data-service/internal/master_sync/service"

	opdClient "master-data-service/internal/opd/client"
	opdRepository "master-data-service/internal/opd/repository"
	opdService "master-data-service/internal/opd/service"

	jabatanClient "master-data-service/internal/jabatan/client"
	jabatanRepository "master-data-service/internal/jabatan/repository"
	jabatanService "master-data-service/internal/jabatan/service"
)

// buildApplication adalah fungsi untuk injeksi repository dan service.
// Tambahkan repository dan service baru di sini.
func buildApplication(
	cfg api.Config,
	logger *slog.Logger,
	db *sql.DB,
) *api.Application {

	// =========================================================
	// REPOSITORY
	// =========================================================

	masterSyncRepo := masterSyncRepository.NewMasterSyncRepository(db, logger)
	dataMasterOpdRepo := opdRepository.NewMasterDataOpdRepository(db, logger)
	dataMasterJabatanRepo := jabatanRepository.NewMasterDataJabatanRepository(db, logger)

	// =========================================================
	// SERVICE
	// =========================================================

	masterSyncSvc := masterSyncService.NewMasterSyncService(
		masterSyncRepo,
		logger,
	)

	simpegClientOpd := opdClient.NewSimpegClient(
		opdClient.SimpegConfig{
			BaseURL: cfg.Simpeg.BaseURL,
			ApiPath: cfg.Simpeg.ApiPath,
		},
		logger,
	)

	simpegClientJabatan := jabatanClient.NewSimpegClient(
		jabatanClient.SimpegConfig{
			BaseURL: cfg.Simpeg.BaseURL,
			ApiPath: cfg.Simpeg.ApiPath,
		},
		logger,
	)

	opdSvc := opdService.NewDataMasterOpdService(
		dataMasterOpdRepo,
		masterSyncSvc,
		simpegClientOpd,
		logger,
	)

	jabatanSvc := jabatanService.NewMasterDataJabatanService(
		dataMasterJabatanRepo,
		masterSyncSvc,
		simpegClientJabatan,
		logger,
	)

	// =========================================================
	// APPLICATION
	// =========================================================

	app := &api.Application{
		Config: cfg,
		Logger: logger,

		MasterSync: masterSyncSvc,

		Opd: opdSvc,

		Jabatan: jabatanSvc,
	}

	return app
}
