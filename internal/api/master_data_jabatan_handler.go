package api

import (
	"net/http"

	jabatanWeb "master-data-service/internal/jabatan/web"
	masterSyncWeb "master-data-service/internal/master_sync/web"
	"master-data-service/internal/web"
)

// FindJabatanByKodeOpd godoc
// @Summary Get Jabatan By Kode OPD
// @Description Get all master Jabatan for a specific OPD.
// @Tags jabatan
// @Produce json
// @Param kodeOpd query string true "Kode OPD"
// @Success 200 {object} web.Response[[]jabatanWeb.DataMasterJabatanResponse]
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /jabatan/find [get]
func (app *Application) FindJabatanByKodeOpd(
	w http.ResponseWriter,
	r *http.Request,
) {
	kodeOpd := r.URL.Query().Get("kodeOpd")

	masterJabatanData, err := app.Jabatan.FindAllJabatanByKodeOpd(
		r.Context(),
		kodeOpd,
	)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	data := web.Response[[]jabatanWeb.DataMasterJabatanResponse]{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "DATA MASTER JABATAN",
		Data:    masterJabatanData,
	}

	if err := app.WriteJSON(w, http.StatusOK, data, nil); err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

// SyncDataMasterJabatan godoc
// @Summary Sync Data Master Jabatan By Kode OPD
// @Description Synchronize master Jabatan data for a specific OPD from the upstream SIMPEG service.
// @Tags jabatan
// @Produce json
// @Param request body masterSyncWeb.MasterSyncRequest true "Sync Pegawai Request"
// @Success 200 {object} web.Response[masterSyncWeb.MasterSyncResponse]
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /jabatan/sync [post]
func (app *Application) SyncDataMasterJabatan(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req masterSyncWeb.MasterSyncRequest

	if err := app.ReadJSON(w, r, &req); err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	syncResult, err := app.Jabatan.SyncDataJabatan(
		r.Context(),
		req,
	)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	data := web.Response[masterSyncWeb.MasterSyncResponse]{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "DATA MASTER JABATAN",
		Data:    syncResult,
	}

	if err := app.WriteJSON(w, http.StatusOK, data, nil); err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}
