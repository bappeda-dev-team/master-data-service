package domain

import "time"

type DataMasterJabatan struct {
	ID uint64

	OpdKode string

	JenisJabatanKode string
	JenisJabatanNama *string

	JabatanKode string
	JabatanNama string

	JenjangJabatanKode *string
	JenjangJabatanNama *string

	SyncID   *string
	SyncFrom *string
	SyncAt   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
