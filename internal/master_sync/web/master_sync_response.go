package web

import (
	"master-data-service/internal/master_sync/domain"
	"time"
)

type MasterSyncResponse struct {
	ID           string            `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	SyncFrom     string            `json:"sync_from" example:"simpeg-adapter"`
	SyncType     domain.SyncType   `json:"sync_type" example:"pegawai"`
	ScopeKey     string            `json:"scope_key" example:"42"`
	StartedAt    time.Time         `json:"started_at" example:"2026-09-01T20:00:00Z"`
	FinishedAt   *time.Time        `json:"finished_at,omitempty" example:"2026-09-01T20:00:13Z"`
	Status       domain.SyncStatus `json:"status" example:"SUCCESS"`
	TotalData    int               `json:"total_data" example:"1250"`
	TotalSuccess int               `json:"total_success" example:"1250"`
	TotalFailed  int               `json:"total_failed" example:"0"`
}
