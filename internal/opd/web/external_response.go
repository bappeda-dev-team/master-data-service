package web

type BaseSimpegResponse struct {
	Title  string              `json:"Title"`
	Status string              `json:"Status"`
	Data   []SimpegOpdResponse `json:"Data"`
}

type SimpegOpdResponse struct {
	OpdKode    string  `json:"opd_kode"`
	OpdLevel   string  `json:"opd_level"`
	OpdNama    string  `json:"opd_nama"`
	OpdJabatan *string `json:"opd_jabatan,omitempty"`
	EselonKode *string `json:"eselon_kode,omitempty"`
	EselonNama *string `json:"eselon_nama,omitempty"`
}
