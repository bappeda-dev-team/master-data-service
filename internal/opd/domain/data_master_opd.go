package domain

import "time"

type DataMasterOpd struct {
	ID         uint64
	OpdKode    string
	OpdLevel   int16
	OpdNama    string
	OpdJabatan *string
	EselonKode *string
	EselonNama *string

	SyncID   *string
	SyncFrom *string
	SyncAt   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
