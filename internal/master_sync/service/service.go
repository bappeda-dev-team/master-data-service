package service

import (
	"context"
	"log/slog"

	"master-data-service/internal/master_sync/domain"
	"master-data-service/internal/master_sync/repository"
)

type MasterSyncService struct {
	Repo   *repository.MasterSyncRepository
	Logger *slog.Logger
}

func NewMasterSyncService(
	repo *repository.MasterSyncRepository,
	// logger taruh paling bawah
	logger *slog.Logger,
) *MasterSyncService {
	return &MasterSyncService{
		Repo:   repo,
		Logger: logger,
	}
}

func (s *MasterSyncService) GetActiveSnapshot(
	ctx context.Context,
	snap domain.DataMasterSync,
) (string, error) {

	snapId, err := s.Repo.FindActiveSnapshot(
		ctx,
		snap,
	)
	if err != nil {
		s.Logger.Error("GetActiveSnapshot")
		return "", err
	}
	if snapId == "" {
		s.Logger.WarnContext(
			ctx,
			"get active snapshot failed",
			"sync_type", snap.SyncType,
			"sync_from", snap.SyncFrom,
			"err", err,
		)
		return "", nil
	}

	return snapId, nil
}
