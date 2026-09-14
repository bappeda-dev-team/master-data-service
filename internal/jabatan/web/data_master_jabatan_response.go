package web

type DataMasterJabatanResponse struct {
	ID uint64 `json:"id"`

	OpdKode string `json:"opd_kode"`

	JenisJabatanKode string  `json:"jenis_jabatan_kode"`
	JenisJabatanNama *string `json:"jenis_jabatan_nama,omitempty"`

	JabatanKode string `json:"jabatan_kode"`
	JabatanNama string `json:"jabatan_nama"`

	JenjangJabatanNama *string `json:"jenjang_jabatan_nama,omitempty"`
}
