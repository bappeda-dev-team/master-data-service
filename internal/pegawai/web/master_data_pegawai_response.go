package web

// import "time"

type DataMasterPegawaiResponse struct {
	ID uint64 `json:"id"`

	PegawaiNIP  string `json:"pegawai_nip"`
	PegawaiNama string `json:"pegawai_nama"`

	PegawaiJenisJabatanKode *string `json:"pegawai_jenis_jabatan_kode,omitempty"`
	PegawaiJabatanKode      *string `json:"pegawai_jabatan_kode,omitempty"`
	PegawaiJabatanTerakhir  *string `json:"pegawai_jabatan_terakhir,omitempty"`

	OpdNama *string `json:"opd_nama,omitempty"`
}
