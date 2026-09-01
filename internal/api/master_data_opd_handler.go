package api

import (
	"net/http"

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
// @Router /healthcheck [get]
func (app *Application) FindAllOpd(w http.ResponseWriter, r *http.Request) {
	masterOpdData, err := app.Opd.FindAllOpd(r.Context())
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}

	data := web.Response[[]opdWeb.DataMasterOpdResponse]{
		Data: masterOpdData,
	}

	err = app.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}
}
