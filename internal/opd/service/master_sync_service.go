package service

import (
	"context"
	"master-data-service/internal/master_sync/domain"
)

type SnapshotService interface {
	GetActiveSnapshot(
		ctx context.Context,
		sync domain.DataMasterSync,
	) (string, error)
}
