package service

import (
	"context"
	"log/slog"

	masterSyncDomain "master-data-service/internal/master_sync/domain"
	"master-data-service/internal/opd/repository"
	// "master-data-service/internal/opd/service/sync"
	"master-data-service/internal/opd/web"
)

type DataMasterOpdService struct {
	Repo        *repository.MasterDataOpdRepository
	SyncService SnapshotService
	// Client       *client.Client
	// SyncExecutor *sync.Registry
	Logger *slog.Logger
}

func NewDataMasterOpdService(
	repo *repository.MasterDataOpdRepository,
	syncService SnapshotService,
	// client *client.Client,
	// syncExecutor *sync.Registry,
	// logger taruh paling bawah
	logger *slog.Logger,
) *DataMasterOpdService {
	return &DataMasterOpdService{
		Repo:        repo,
		SyncService: syncService,
		// Client:       client,
		// SyncExecutor: syncExecutor,
		Logger: logger,
	}
}

func (s *DataMasterOpdService) FindAllOpd(
	ctx context.Context,
) ([]web.DataMasterOpdResponse, error) {
	s.Logger.Info("FindAllOpd")

	const syncFrom = "SIMPEG"

	snapshot := masterSyncDomain.DataMasterSync{
		SyncType: masterSyncDomain.SyncTypeOpd,
		SyncFrom: syncFrom,
		Status:   masterSyncDomain.SyncStatusSuccess,
	}

	// FIND LATEST SNAPSHOT
	activeSnapshotId, err := s.SyncService.GetActiveSnapshot(ctx, snapshot)
	if err != nil {
		return nil, err
	}
	masterOpd, err := s.Repo.FindAllOpdBySnapshotId(
		ctx,
		activeSnapshotId,
	)
	if err != nil {
		return nil, err
	}
	response := make([]web.DataMasterOpdResponse, 0, len(masterOpd))
	for _, opd := range masterOpd {
		response = append(response,
			web.DataMasterOpdResponse{
				ID:       opd.ID,
				OpdKode:  opd.OpdKode,
				OpdLevel: opd.OpdLevel,
				OpdNama:  opd.OpdNama,
			},
		)
	}
	return response, nil

}
