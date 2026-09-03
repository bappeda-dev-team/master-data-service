package web

import "time"

type DataMasterPegawaiResponse struct {
	ID uint64 `json:"id"`

	OpdKode string `json:"opd_kode"`

	PegawaiNIP  string `json:"pegawai_nip"`
	PegawaiNama string `json:"pegawai_nama"`

	PegawaiJenisKelamin *string `json:"pegawai_jenis_kelamin,omitempty"`

	PegawaiLahirTempat  *string    `json:"pegawai_lahir_tempat,omitempty"`
	PegawaiLahirTanggal *time.Time `json:"pegawai_lahir_tanggal,omitempty"`

	PegawaiPensiunTMT *time.Time `json:"pegawai_pensiun_tmt,omitempty"`

	PegawaiAktifID *string `json:"pegawai_aktif_id,omitempty"`
	PegawaiStatus  *string `json:"pegawai_status,omitempty"`
	PegawaiJenis   *string `json:"pegawai_jenis,omitempty"`

	PegawaiGolonganKode    *string    `json:"pegawai_golongan_kode,omitempty"`
	PegawaiGolonganNama    *string    `json:"pegawai_golongan_nama,omitempty"`
	PegawaiGolonganPangkat *string    `json:"pegawai_golongan_pangkat,omitempty"`
	PegawaiGolonganTMT     *time.Time `json:"pegawai_golongan_tmt,omitempty"`

	PegawaiJenisJabatanKode *string `json:"pegawai_jenis_jabatan_kode,omitempty"`
	PegawaiJenisJabatanNama *string `json:"pegawai_jenis_jabatan_nama,omitempty"`

	PegawaiJabatanKode     *string `json:"pegawai_jabatan_kode,omitempty"`
	PegawaiJabatanTerakhir *string `json:"pegawai_jabatan_terakhir,omitempty"`

	PegawaiEselonKode *string `json:"pegawai_eselon_kode,omitempty"`
	PegawaiEselonNama *string `json:"pegawai_eselon_nama,omitempty"`

	OpdLevel *int16  `json:"opd_level,omitempty"`
	OpdNama  *string `json:"opd_nama,omitempty"`
}
