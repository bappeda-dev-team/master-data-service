package api

import (
	"net/http"

	masterSyncWeb "master-data-service/internal/master_sync/web"
	pegawaiWeb "master-data-service/internal/pegawai/web"
	"master-data-service/internal/web"
)

// FindPegawaiByKodeOpd godoc
// @Summary Get all Pegawai
// @Description Get all Pegawai By Kode Opd
// @Tags pegawai
// @Produce json
// @Param kodeOpd query string true "Kode OPD"
// @Success 200	{object} web.Response[[]pegawaiWeb.DataMasterPegawaiResponse]
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /pegawai/find [get]
func (app *Application) FindPegawaiByKodeOpd(w http.ResponseWriter, r *http.Request) {
	kodeOpd := r.URL.Query().Get("kodeOpd")

	masterDataPegawai, err := app.Pegawai.FindAllPegawaiByKodeOpd(
		r.Context(),
		kodeOpd,
	)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	data := web.Response[[]pegawaiWeb.DataMasterPegawaiResponse]{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "DATA MASTER PEGAWAI",
		Data:    masterDataPegawai,
	}

	err = app.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

// SyncDataMasterPegawai godoc
// @Summary Sync Data Master Pegawai
// @Description Synchronize master Pegawai data for a specific OPD from the upstream SIMPEG service.
// @Tags pegawai
// @Produce json
// @Param kodeOpd query string true "Kode OPD"
// @Success 200	{object} web.Response[masterSyncWeb.MasterSyncResponse]
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /pegawai/sync [post]
func (app *Application) SyncDataMasterPegawai(w http.ResponseWriter, r *http.Request) {
	kodeOpd := r.URL.Query().Get("kodeOpd")

	syncResult, err := app.Pegawai.SyncDataPegawai(
		r.Context(),
		kodeOpd,
	)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	data := web.Response[masterSyncWeb.MasterSyncResponse]{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "SYNC DATA MASTER PEGAWAI RESULT",
		Data:    syncResult,
	}

	err = app.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}
