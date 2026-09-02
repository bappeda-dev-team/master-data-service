package web

type BaseSimpegResponse struct {
	Title  string                  `json:"Title"`
	Status string                  `json:"Status"`
	Data   []SimpegJabatanResponse `json:"Data"`
}

type SimpegJabatanResponse struct {
	OpdKode            string `json:"opd_kode"`
	OpdNama            string `json:"opd_nama"`
	JenisJabatanKode   string `json:"jenis_jabatan_kode"`
	JenisJabatanNama   string `json:"jenis_jabatan_nama"`
	JabatanKode        string `json:"jabatan_kode"`
	JabatanNama        string `json:"jabatan_nama"`
	JenjangJabatanKode string `json:"jenjang_jabatan_kode"`
	JenjangJabatanNama string `json:"jenjang_jabatan_nama"`
}
