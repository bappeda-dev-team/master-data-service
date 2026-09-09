package web

type MappingOpdResponse struct {
	ID         int64  `json:"id"`
	Sumber     string `json:"sumber"`
	KodeSumber string `json:"kode_sumber"`
	KodeMaster string `json:"kode_master"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type CreateMappingOpdRequest struct {
	Sumber     string `json:"sumber"`
	KodeSumber string `json:"kode_sumber"`
	KodeMaster string `json:"kode_master"`
}

type UpdateMappingOpdRequest struct {
	Sumber     string `json:"sumber"`
	KodeSumber string `json:"kode_sumber"`
	KodeMaster string `json:"kode_master"`
}
