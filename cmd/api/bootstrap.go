package main

import (
	"database/sql"
	"log/slog"

	"master-data-service/internal/api"

	masterSyncRepository "master-data-service/internal/master_sync/repository"
	masterSyncService "master-data-service/internal/master_sync/service"

	opdRepository "master-data-service/internal/opd/repository"
	opdService "master-data-service/internal/opd/service"
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

	// =========================================================
	// SERVICE
	// =========================================================

	masterSyncSvc := masterSyncService.NewMasterSyncService(
		masterSyncRepo,
		logger,
	)

	opdSvc := opdService.NewDataMasterOpdService(
		dataMasterOpdRepo,
		masterSyncSvc,
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
	}

	return app
}
