package web

type BaseSimpegResponse struct {
	Title  string                  `json:"Title"`
	Status string                  `json:"Status"`
	Data   []SimpegPegawaiResponse `json:"Data"`
}

type SimpegPegawaiResponse struct {
	PegawaiNIP  string `json:"pegawai_nip"`
	PegawaiNama string `json:"pegawai_nama"`

	PegawaiJenisKelamin string `json:"pegawai_jenis_kelamin"`

	PegawaiLahirTempat  string `json:"pegawai_lahir_tempat"`
	PegawaiLahirTanggal string `json:"pegawai_lahir_tanggal"`
	PegawaiPensiunTMT   string `json:"pegawai_pensiun_tmt"`

	PegawaiAktifID string `json:"pegawai_aktif_id"`
	PegawaiStatus  string `json:"pegawai_status"`
	PegawaiJenis   string `json:"pegawai_jenis"`

	PegawaiGolonganKode    string `json:"pegawai_golongan_kode"`
	PegawaiGolonganNama    string `json:"pegawai_golongan_nama"`
	PegawaiGolonganPangkat string `json:"pegawai_golongan_pangkat"`
	PegawaiGolonganTMT     string `json:"pegawai_golongan_tmt"`

	PegawaiJenisJabatanKode string `json:"pegawai_jenis_jabatan_kode"`
	PegawaiJenisJabatanNama string `json:"pegawai_jenis_jabatan_nama"`

	PegawaiJabatanKode     string `json:"pegawai_jabatan_kode"`
	PegawaiJabatanTerakhir string `json:"pegawai_jabatan_terakhir"`

	PegawaiEselonKode string `json:"pegawai_eselon_kode"`
	PegawaiEselonNama string `json:"pegawai_eselon_nama"`

	OpdKode  string `json:"opd_kode"`
	OpdLevel string `json:"opd_level"`
	OpdNama  string `json:"opd_nama"`
}
