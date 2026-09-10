package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"master-data-service/internal/pegawai/client"
	"master-data-service/internal/pegawai/domain"
	"master-data-service/internal/pegawai/repository"
	"master-data-service/internal/pegawai/web"
	"strconv"
	"time"

	masterSyncDomain "master-data-service/internal/master_sync/domain"
	masterSyncWeb "master-data-service/internal/master_sync/web"

	"github.com/google/uuid"
)

type MasterDataPegawaiService struct {
	Repo        *repository.MasterDataPegawaiRepository
	SyncService SnapshotService
	Client      *client.SimpegClient
	Logger      *slog.Logger
}

func NewMasterDataPegawaiService(
	repo *repository.MasterDataPegawaiRepository,
	syncService SnapshotService,
	simpegClient *client.SimpegClient,
	logger *slog.Logger,
) *MasterDataPegawaiService {
	return &MasterDataPegawaiService{
		Repo:        repo,
		SyncService: syncService,
		Client:      simpegClient,
		Logger:      logger,
	}
}

func (s *MasterDataPegawaiService) SyncDataPegawai(
	ctx context.Context,
	kodeOpd string,
) (masterSyncWeb.MasterSyncResponse, error) {
	s.Logger.Info("SyncDataPegawai")

	const syncFrom = "SIMPEG"

	// CREATE SYNC RECORD
	syncID := uuid.New().String()
	sync := masterSyncDomain.DataMasterSync{
		ID:        syncID,
		SyncFrom:  syncFrom,
		SyncType:  masterSyncDomain.SyncTypePegawai,
		ScopeKey:  kodeOpd,
		StartedAt: time.Now(),
		Status:    masterSyncDomain.SyncStatusRunning,
	}

	if err := s.SyncService.Create(ctx, sync); err != nil {
		return masterSyncWeb.MasterSyncResponse{}, err
	}

	simpegPegawais, err := s.Client.MasterPegawaiSimpeg(ctx, kodeOpd)
	if err != nil {
		s.finishSync(
			ctx,
			sync,
			0,
			0,
			len(simpegPegawais),
			masterSyncDomain.SyncStatusFailed,
		)

		return masterSyncWeb.MasterSyncResponse{}, err
	}

	pegawais := make([]domain.DataMasterPegawai, 0, len(simpegPegawais))

	for _, peg := range simpegPegawais {
		pegawaiLahirTanggal, err := parseDate(peg.PegawaiLahirTanggal)
		if err != nil {
			s.finishSync(
				ctx,
				sync,
				0,
				0,
				len(simpegPegawais),
				masterSyncDomain.SyncStatusFailed,
			)

			return masterSyncWeb.MasterSyncResponse{}, fmt.Errorf(
				"parse pegawai_lahir_tanggal for NIP %s: %w",
				peg.PegawaiNIP,
				err,
			)
		}

		pegawaiPensiunTMT, err := parseDate(peg.PegawaiPensiunTMT)
		if err != nil {
			s.finishSync(
				ctx,
				sync,
				0,
				0,
				len(simpegPegawais),
				masterSyncDomain.SyncStatusFailed,
			)

			return masterSyncWeb.MasterSyncResponse{}, fmt.Errorf(
				"parse pegawai_pensiun_tmt for NIP %s: %w",
				peg.PegawaiNIP,
				err,
			)
		}

		pegawaiGolonganTMT, err := parseDate(peg.PegawaiGolonganTMT)
		if err != nil {
			s.finishSync(
				ctx,
				sync,
				0,
				0,
				len(simpegPegawais),
				masterSyncDomain.SyncStatusFailed,
			)

			return masterSyncWeb.MasterSyncResponse{}, fmt.Errorf(
				"parse pegawai_golongan_tmt for NIP %s: %w",
				peg.PegawaiNIP,
				err,
			)
		}

		opdLevel, err := parseInt16(peg.OpdLevel)
		if err != nil {
			s.finishSync(
				ctx,
				sync,
				0,
				0,
				len(simpegPegawais),
				masterSyncDomain.SyncStatusFailed,
			)

			return masterSyncWeb.MasterSyncResponse{}, fmt.Errorf(
				"parse opd_level for NIP %s: %w",
				peg.PegawaiNIP,
				err,
			)
		}

		pegawais = append(pegawais, domain.DataMasterPegawai{
			OpdKode: peg.OpdKode,

			PegawaiNIP:  peg.PegawaiNIP,
			PegawaiNama: peg.PegawaiNama,

			PegawaiJenisKelamin: stringPtr(peg.PegawaiJenisKelamin),

			PegawaiLahirTempat:  stringPtr(peg.PegawaiLahirTempat),
			PegawaiLahirTanggal: pegawaiLahirTanggal,

			PegawaiPensiunTMT: pegawaiPensiunTMT,

			PegawaiAktifID: stringPtr(peg.PegawaiAktifID),
			PegawaiStatus:  stringPtr(peg.PegawaiStatus),
			PegawaiJenis:   stringPtr(peg.PegawaiJenis),

			PegawaiGolonganKode:    stringPtr(peg.PegawaiGolonganKode),
			PegawaiGolonganNama:    stringPtr(peg.PegawaiGolonganNama),
			PegawaiGolonganPangkat: stringPtr(peg.PegawaiGolonganPangkat),
			PegawaiGolonganTMT:     pegawaiGolonganTMT,

			PegawaiJenisJabatanKode: stringPtr(peg.PegawaiJenisJabatanKode),
			PegawaiJenisJabatanNama: stringPtr(peg.PegawaiJenisJabatanNama),

			PegawaiJabatanKode:     stringPtr(peg.PegawaiJabatanKode),
			PegawaiJabatanTerakhir: stringPtr(peg.PegawaiJabatanTerakhir),

			PegawaiEselonKode: stringPtr(peg.PegawaiEselonKode),
			PegawaiEselonNama: stringPtr(peg.PegawaiEselonNama),

			OpdLevel: opdLevel,
			OpdNama:  stringPtr(peg.OpdNama),
		})
	}

	// INSERT DATA DARI SIMPEG RESPONSE
	inserted, err := s.Repo.BulkInsert(ctx, pegawais, syncID, syncFrom)
	if err != nil {
		s.finishSync(ctx, sync, len(pegawais), 0, len(pegawais), masterSyncDomain.SyncStatusFailed)
		return masterSyncWeb.MasterSyncResponse{}, err
	}

	// FINISH SYNC
	if err := s.finishSync(ctx, sync, len(pegawais), inserted, 0, masterSyncDomain.SyncStatusSuccess); err != nil {
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
		TotalData:    len(pegawais),
		TotalSuccess: inserted,
		TotalFailed:  0,
	}, nil

}

// finishSync mengupdate status data_master_sync
func (s *MasterDataPegawaiService) finishSync(
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

func parseDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}

	t, err := time.Parse("02-01-2006", value)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func parseInt16(value string) (*int16, error) {
	if value == "" {
		return nil, nil
	}

	n, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return nil, err
	}

	result := int16(n)

	return &result, nil
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

func (s *MasterDataPegawaiService) FindAllPegawaiByKodeOpd(
	ctx context.Context,
	kodeOpd string,
) ([]web.DataMasterPegawaiResponse, error) {
	const op = "pegawai-service-FindAllPegawaiByKodeOpd"
	const syncFrom = "SIMPEG"

	s.Logger.Info(
		"find all pegawai by kode opd",
		"op", op,
		"opd_kode", kodeOpd,
	)

	snapshot := masterSyncDomain.DataMasterSync{
		SyncType: masterSyncDomain.SyncTypePegawai,
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

			return []web.DataMasterPegawaiResponse{}, nil
		}

		s.Logger.Error(
			"failed to find active snapshot",
			"op", op,
			"sync_type", snapshot.SyncType,
			"sync_from", snapshot.SyncFrom,
			"scope_key", snapshot.ScopeKey,
			"error", err,
		)

		return nil, err
	}

	// FIND PEGAWAI BY OPD KODE
	pegawais, err := s.Repo.FindAllPegawaiBySnapshotIdAndKodeOpd(
		ctx,
		snapshotID,
		kodeOpd,
	)
	if err != nil {
		s.Logger.Error(
			"failed to find pegawai",
			"op", op,
			"snapshot_id", snapshotID,
			"opd_kode", kodeOpd,
			"error", err,
		)

		return nil, err
	}

	response := make(
		[]web.DataMasterPegawaiResponse,
		0,
		len(pegawais),
	)

	for _, pegawai := range pegawais {
		response = append(
			response,
			web.DataMasterPegawaiResponse{
				ID:                      pegawai.ID,
				OpdKode:                 pegawai.OpdKode,
				PegawaiNIP:              pegawai.PegawaiNIP,
				PegawaiNama:             pegawai.PegawaiNama,
				PegawaiJenisKelamin:     pegawai.PegawaiJenisKelamin,
				PegawaiLahirTempat:      pegawai.PegawaiLahirTempat,
				PegawaiLahirTanggal:     pegawai.PegawaiLahirTanggal,
				PegawaiPensiunTMT:       pegawai.PegawaiPensiunTMT,
				PegawaiAktifID:          pegawai.PegawaiAktifID,
				PegawaiStatus:           pegawai.PegawaiStatus,
				PegawaiJenis:            pegawai.PegawaiJenis,
				PegawaiGolonganKode:     pegawai.PegawaiGolonganKode,
				PegawaiGolonganNama:     pegawai.PegawaiGolonganNama,
				PegawaiGolonganPangkat:  pegawai.PegawaiGolonganPangkat,
				PegawaiGolonganTMT:      pegawai.PegawaiGolonganTMT,
				PegawaiJenisJabatanKode: pegawai.PegawaiJenisJabatanKode,
				PegawaiJenisJabatanNama: pegawai.PegawaiJenisJabatanNama,
				PegawaiJabatanKode:      pegawai.PegawaiJabatanKode,
				PegawaiJabatanTerakhir:  pegawai.PegawaiJabatanTerakhir,
				PegawaiEselonKode:       pegawai.PegawaiEselonKode,
				PegawaiEselonNama:       pegawai.PegawaiEselonNama,
				OpdLevel:                pegawai.OpdLevel,
				OpdNama:                 pegawai.OpdNama,
			},
		)
	}

	s.Logger.Debug(
		"successfully found pegawai",
		"op", op,
		"snapshot_id", snapshotID,
		"opd_kode", kodeOpd,
		"count", len(response),
	)

	return response, nil
}
