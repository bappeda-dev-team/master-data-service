package web

type MasterSyncRequest struct {
	SyncFrom string `json:"sync_from" example:"SIMPEG"`
	KodeOpd  string `json:"kode_opd" example:"5.10.23"`
}
