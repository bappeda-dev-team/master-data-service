package domain

import "time"

type DataMasterPegawai struct {
	ID uint64

	OpdKode string

	PegawaiNIP  string
	PegawaiNama string

	PegawaiJenisKelamin *string

	PegawaiLahirTempat  *string
	PegawaiLahirTanggal *time.Time

	PegawaiPensiunTMT *time.Time

	PegawaiAktifID *string
	PegawaiStatus  *string
	PegawaiJenis   *string

	PegawaiGolonganKode    *string
	PegawaiGolonganNama    *string
	PegawaiGolonganPangkat *string
	PegawaiGolonganTMT     *time.Time

	PegawaiJenisJabatanKode *string
	PegawaiJenisJabatanNama *string

	PegawaiJabatanKode     *string
	PegawaiJabatanTerakhir *string

	PegawaiEselonKode *string
	PegawaiEselonNama *string

	OpdLevel *int16
	OpdNama  *string

	SyncID   *string
	SyncFrom *string
	SyncAt   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
