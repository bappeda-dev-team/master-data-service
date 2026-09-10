package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"master-data-service/internal/jabatan/client"
	"master-data-service/internal/jabatan/domain"
	"master-data-service/internal/jabatan/repository"
	masterSyncDomain "master-data-service/internal/master_sync/domain"
	masterSyncWeb "master-data-service/internal/master_sync/web"

	// "master-data-service/internal/opd/service/sync"
	"master-data-service/internal/jabatan/web"
)

type MasterDataJabatanService struct {
	Repo        *repository.MasterDataJabatanRepository
	SyncService SnapshotService
	Client      *client.SimpegClient
	// SyncExecutor *sync.Registry
	Logger *slog.Logger
}

func NewMasterDataJabatanService(
	repo *repository.MasterDataJabatanRepository,
	syncService SnapshotService,
	simpegClient *client.SimpegClient,
	// syncExecutor *sync.Registry,
	// logger taruh paling bawah
	logger *slog.Logger,
) *MasterDataJabatanService {
	return &MasterDataJabatanService{
		Repo:        repo,
		SyncService: syncService,
		Client:      simpegClient,
		// SyncExecutor:  syncExecutor,
		Logger: logger,
	}
}

func (s *MasterDataJabatanService) SyncDataJabatan(
	ctx context.Context,
	kodeOpd string,
) (masterSyncWeb.MasterSyncResponse, error) {
	s.Logger.Info("SyncDataJabatan")

	const syncFrom = "SIMPEG"

	// CREATE SYNC RECORD
	syncID := uuid.New().String()
	sync := masterSyncDomain.DataMasterSync{
		ID:        syncID,
		SyncFrom:  syncFrom,
		SyncType:  masterSyncDomain.SyncTypeJabatan,
		ScopeKey:  kodeOpd,
		StartedAt: time.Now(),
		Status:    masterSyncDomain.SyncStatusRunning,
	}

	if err := s.SyncService.Create(ctx, sync); err != nil {
		return masterSyncWeb.MasterSyncResponse{}, err
	}

	// CALL CLIENT UNTUK MENGAMBIL DATA OPD DARI SIMPEG
	// kodeOpd = kode opd dari master data opd simpeg
	simpegJabatans, err := s.Client.MasterJabatanSimpeg(ctx, kodeOpd)
	if err != nil {
		s.finishSync(ctx, sync, 0, 0, len(simpegJabatans), masterSyncDomain.SyncStatusFailed)
		return masterSyncWeb.MasterSyncResponse{}, err
	}

	// MAP WEB RESPONSE -> DOMAIN
	jabatans := make([]domain.DataMasterJabatan, 0, len(simpegJabatans))

	for _, jb := range simpegJabatans {
		jabatans = append(jabatans, domain.DataMasterJabatan{
			OpdKode:            jb.OpdKode,
			JenisJabatanKode:   jb.JenisJabatanKode,
			JenisJabatanNama:   stringPtr(jb.JenisJabatanNama),
			JabatanKode:        jb.JabatanKode,
			JabatanNama:        jb.JabatanNama,
			JenjangJabatanKode: stringPtr(jb.JenjangJabatanKode),
			JenjangJabatanNama: stringPtr(jb.JenjangJabatanNama),
		})
	}

	// INSERT DATA DARI SIMPEG RESPONSE
	inserted, err := s.Repo.BulkInsert(ctx, jabatans, syncID, syncFrom)
	if err != nil {
		s.finishSync(ctx, sync, len(jabatans), 0, len(jabatans), masterSyncDomain.SyncStatusFailed)
		return masterSyncWeb.MasterSyncResponse{}, err
	}

	// FINISH SYNC
	if err := s.finishSync(ctx, sync, len(jabatans), inserted, 0, masterSyncDomain.SyncStatusSuccess); err != nil {
		return masterSyncWeb.MasterSyncResponse{}, err
	}

	return masterSyncWeb.MasterSyncResponse{
		ID:           syncID,
		SyncFrom:     syncFrom,
		SyncType:     sync.SyncType,
		ScopeKey:     sync.ScopeKey,
		StartedAt:    sync.StartedAt,
		FinishedAt:   sync.FinishedAt,
		Status:       masterSyncDomain.SyncStatusSuccess,
		TotalData:    len(jabatans),
		TotalSuccess: inserted,
		TotalFailed:  0,
	}, nil
}

// finishSync mengupdate status data_master_sync
func (s *MasterDataJabatanService) finishSync(
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

func (s *MasterDataJabatanService) FindAllJabatanByKodeOpd(
	ctx context.Context,
	kodeOpd string,
) ([]web.DataMasterJabatanResponse, error) {
	const op = "jabatan-service-FindAllJabatanByKodeOpd"
	const syncFrom = "SIMPEG"

	s.Logger.Info(
		"find all jabatan by kode opd",
		"op", op,
		"opd_kode", kodeOpd,
	)

	snapshot := masterSyncDomain.DataMasterSync{
		SyncType: masterSyncDomain.SyncTypeJabatan,
		SyncFrom: syncFrom,
		ScopeKey: kodeOpd,
		Status:   masterSyncDomain.SyncStatusSuccess,
	}

	// FIND LATEST SNAPSHOT
	snapshotID, err := s.SyncService.GetActiveSnapshot(ctx, snapshot)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.Logger.Debug(
				"no active snapshot found",
				"op", op,
				"sync_type", snapshot.SyncType,
				"sync_from", snapshot.SyncFrom,
				"scope_key", snapshot.ScopeKey,
			)

			return []web.DataMasterJabatanResponse{}, nil
		}

		s.Logger.Error(
			"failed to find active snapshot",
			"op", op,
			"sync_type", snapshot.SyncType,
			"sync_from", snapshot.SyncFrom,
			"error", err,
		)

		return nil, err
	}

	// FIND JABATAN BY OPD KODE
	jabatans, err := s.Repo.FindAllJabatanBySnapshotIdAndKodeOpd(
		ctx,
		snapshotID,
		kodeOpd,
	)
	if err != nil {
		s.Logger.Error(
			"failed to find jabatan",
			"op", op,
			"snapshot_id", snapshotID,
			"opd_kode", kodeOpd,
			"error", err,
		)

		return nil, err
	}

	response := make(
		[]web.DataMasterJabatanResponse,
		0,
		len(jabatans),
	)

	for _, jabatan := range jabatans {
		response = append(
			response,
			web.DataMasterJabatanResponse{
				ID:                 jabatan.ID,
				OpdKode:            jabatan.OpdKode,
				JenisJabatanKode:   jabatan.JenisJabatanKode,
				JenisJabatanNama:   jabatan.JenisJabatanNama,
				JabatanKode:        jabatan.JabatanKode,
				JabatanNama:        jabatan.JabatanNama,
				JenjangJabatanKode: jabatan.JenjangJabatanKode,
				JenjangJabatanNama: jabatan.JenjangJabatanNama,
			},
		)
	}

	s.Logger.Debug(
		"successfully found jabatan",
		"op", op,
		"snapshot_id", snapshotID,
		"opd_kode", kodeOpd,
		"count", len(response),
	)

	return response, nil
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
