package web

type DataMasterOpdResponse struct {
	ID              uint64  `json:"id"`
	OpdKode         string  `json:"opd_kode"`
	OpdLevel        int16   `json:"opd_level"`
	OpdNama         string  `json:"opd_nama"`
	OpdJabatan      *string `json:"opd_jabatan,omitempty"`
	EselonKode      *string `json:"eselon_kode,omitempty"`
	EselonNama      *string `json:"eselon_nama,omitempty"`
	KodeNomenklatur *string `json:"kode_nomenklatur"`
}
