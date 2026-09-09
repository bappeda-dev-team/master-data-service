package api

import (
	"net/http"

	masterSyncWeb "master-data-service/internal/master_sync/web"
	opdWeb "master-data-service/internal/opd/web"
	"master-data-service/internal/web"
)

// FindAllOpd godoc
// @Summary Get all Opd
// @Description Get all opd
// @Tags opd
// @Produce json
// @Success 200	{object} web.Response[[]opdWeb.DataMasterOpdResponse]
// @Failure 500 {object} web.ErrorResponse
// @Router /opd/findall [get]
func (app *Application) FindAllOpd(w http.ResponseWriter, r *http.Request) {
	masterOpdData, err := app.Opd.FindAllOpd(r.Context())
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	data := web.Response[[]opdWeb.DataMasterOpdResponse]{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "DATA MASTER OPD",
		Data:    masterOpdData,
	}

	err = app.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

// SyncDataMasterOpd godoc
// @Summary Sync Data Master Opd
// @Description Update OPD Data From Upstream
// @Tags opd
// @Produce json
// @Success 200	{object} web.Response[masterSyncWeb.MasterSyncResponse]
// @Failure 500 {object} web.ErrorResponse
// @Router /opd/sync [post]
func (app *Application) SyncDataMasterOpd(w http.ResponseWriter, r *http.Request) {
	syncResult, err := app.Opd.SyncDataOpd(r.Context())
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	data := web.Response[masterSyncWeb.MasterSyncResponse]{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: SyncMessage(syncResult.Status),
		Data:    syncResult,
	}

	err = app.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}
