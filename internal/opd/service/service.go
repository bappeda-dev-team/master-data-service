package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"

	masterSyncDomain "master-data-service/internal/master_sync/domain"
	masterSyncWeb "master-data-service/internal/master_sync/web"
	"master-data-service/internal/opd/client"
	"master-data-service/internal/opd/domain"
	"master-data-service/internal/opd/repository"

	// "master-data-service/internal/opd/service/sync"
	"master-data-service/internal/opd/web"
)

type DataMasterOpdService struct {
	Repo        *repository.MasterDataOpdRepository
	SyncService SnapshotService
	Client      *client.SimpegClient
	// SyncExecutor *sync.Registry
	Logger *slog.Logger
}

func NewDataMasterOpdService(
	repo *repository.MasterDataOpdRepository,
	syncService SnapshotService,
	simpegClient *client.SimpegClient,
	// syncExecutor *sync.Registry,
	// logger taruh paling bawah
	logger *slog.Logger,
) *DataMasterOpdService {
	return &DataMasterOpdService{
		Repo:        repo,
		SyncService: syncService,
		Client:      simpegClient,
		// SyncExecutor:  syncExecutor,
		Logger: logger,
	}
}

func (s *DataMasterOpdService) SyncDataOpd(
	ctx context.Context,
) (masterSyncWeb.MasterSyncResponse, error) {
	s.Logger.Info("SyncDataOpd")

	const syncFrom = "SIMPEG"

	// CREATE SYNC RECORD
	syncID := uuid.New().String()
	sync := masterSyncDomain.DataMasterSync{
		ID:        syncID,
		SyncFrom:  syncFrom,
		SyncType:  masterSyncDomain.SyncTypeOpd,
		StartedAt: time.Now(),
		Status:    masterSyncDomain.SyncStatusRunning,
	}

	if err := s.SyncService.Create(ctx, sync); err != nil {
		return masterSyncWeb.MasterSyncResponse{}, err
	}

	// CALL CLIENT UNTUK MENGAMBIL DATA OPD DARI SIMPEG
	simpegOpds, err := s.Client.MasterOpdSimpeg(ctx)
	if err != nil {
		s.finishSync(ctx, sync, 0, 0, len(simpegOpds), masterSyncDomain.SyncStatusFailed)
		return masterSyncWeb.MasterSyncResponse{}, err
	}

	// MAP WEB RESPONSE -> DOMAIN
	opds := make([]domain.DataMasterOpd, 0, len(simpegOpds))
	for _, so := range simpegOpds {
		opdLevelInt, errLevel := strconv.Atoi(so.OpdLevel)
		if errLevel != nil {
			return masterSyncWeb.MasterSyncResponse{}, err
		}
		opds = append(opds, domain.DataMasterOpd{
			OpdKode:    so.OpdKode,
			OpdLevel:   int16(opdLevelInt),
			OpdNama:    so.OpdNama,
			OpdJabatan: so.OpdJabatan,
			EselonKode: so.EselonKode,
			EselonNama: so.EselonNama,
			SyncID:     &syncID,
		})
	}

	totalData := len(opds)

	// INSERT DATA DARI SIMPEG RESPONSE
	inserted, err := s.Repo.BulkInsert(ctx, opds, syncID, syncFrom)
	if err != nil {
		s.finishSync(ctx, sync, len(opds), 0, len(opds), masterSyncDomain.SyncStatusFailed)
		return masterSyncWeb.MasterSyncResponse{}, err
	}

	status := masterSyncDomain.SyncStatusSuccess
	totalFailed := totalData - inserted

	if totalFailed > 0 {
		status = masterSyncDomain.SyncStatusPartial
	}

	// FINISH SYNC
	if err := s.finishSync(ctx, sync, len(opds), inserted, 0, masterSyncDomain.SyncStatusSuccess); err != nil {
		return masterSyncWeb.MasterSyncResponse{}, err
	}

	return masterSyncWeb.MasterSyncResponse{
		ID:           syncID,
		SyncFrom:     syncFrom,
		SyncType:     sync.SyncType,
		StartedAt:    sync.StartedAt,
		FinishedAt:   sync.FinishedAt,
		Status:       status,
		TotalData:    totalData,
		TotalSuccess: inserted,
		TotalFailed:  totalFailed,
	}, nil
}

// finishSync mengupdate status data_master_sync
func (s *DataMasterOpdService) finishSync(
	ctx context.Context,
	sync masterSyncDomain.DataMasterSync,
	totalData int,
	totalSuccess int,
	totalFailed int,
	status masterSyncDomain.SyncStatus,
) error {
	now := time.Now()
	sync.FinishedAt = &now
	sync.Status = status
	sync.TotalData = totalData
	sync.TotalSuccess = totalSuccess
	sync.TotalFailed = totalFailed

	return s.SyncService.Finish(ctx, sync)
}

// ONLY GET OPD LEVEL 1
func (s *DataMasterOpdService) FindAllOpd(
	ctx context.Context,
) ([]web.DataMasterOpdResponse, error) {
	s.Logger.Info("FindAllOpd LEVEL 1")

	const syncFrom = "SIMPEG"

	snapshot := masterSyncDomain.DataMasterSync{
		SyncType: masterSyncDomain.SyncTypeOpd,
		SyncFrom: syncFrom,
		Status:   masterSyncDomain.SyncStatusSuccess,
	}

	// FIND LATEST SNAPSHOT
	activeSnapshotId, err := s.SyncService.GetActiveSnapshot(ctx, snapshot)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.Debug(
				"no active snapshot found",
				"sync_type", snapshot.SyncType,
				"sync_from", snapshot.SyncFrom,
			)

			return []web.DataMasterOpdResponse{}, nil
		}

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
