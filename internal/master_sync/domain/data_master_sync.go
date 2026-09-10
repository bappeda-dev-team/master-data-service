package domain

import "time"

type DataMasterSync struct {
	ID           string
	SyncFrom     string
	SyncType     SyncType
	ScopeKey     string
	StartedAt    time.Time
	FinishedAt   *time.Time
	Status       SyncStatus
	TotalData    int
	TotalSuccess int
	TotalFailed  int
}
