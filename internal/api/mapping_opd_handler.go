package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"master-data-service/internal/opd/domain"
	mappingOpdService "master-data-service/internal/opd/service"
	opdWeb "master-data-service/internal/opd/web"
	"master-data-service/internal/web"
)

// FindAllMappingOpd godoc
// @Summary Get all mapping opd
// @Description Get all mapping opd
// @Tags mapping-opd
// @Produce json
// @Success 200 {object} web.Response[[]opdWeb.MappingOpdResponse]
// @Failure 500 {object} web.ErrorResponse
// @Router /mapping-opd/findall [get]
func (app *Application) FindAllMappingOpd(
	w http.ResponseWriter,
	r *http.Request,
) {
	mappings, err := app.MappingOpd.FindAll(r.Context())
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	data := web.Response[[]opdWeb.MappingOpdResponse]{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "DATA MAPPING OPD",
		Data:    mappings,
	}

	if err := app.WriteJSON(w, http.StatusOK, data, nil); err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

// FindMappingOpdByKodeSumber godoc
// @Summary Get mapping opd by kode sumber
// @Description Get mapping opd by kode sumber
// @Tags mapping-opd
// @Produce json
// @Param sumber query string true "Sumber mapping"
// @Param kode_sumber query string true "Kode OPD dari sumber"
// @Success 200 {object} web.Response[opdWeb.MappingOpdResponse]
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /mapping-opd/find-by-kode-sumber [get]
func (app *Application) FindMappingOpdByKodeSumber(
	w http.ResponseWriter,
	r *http.Request,
) {
	q := r.URL.Query()

	sumber := q.Get("sumber")
	kodeSumber := q.Get("kode_sumber")

	mapping, err := app.MappingOpd.FindBySumberAndKodeSumber(
		r.Context(),
		sumber,
		kodeSumber,
	)
	if err != nil {
		if errors.Is(err, mappingOpdService.ErrMappingOpdNotFound) {
			app.NotFoundResponse(w, r)
			return
		}

		app.ServerErrorResponse(w, r, err)
		return
	}

	data := web.Response[*opdWeb.MappingOpdResponse]{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "DATA MAPPING OPD",
		Data:    mapping,
	}

	if err := app.WriteJSON(w, http.StatusOK, data, nil); err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

// FindMappingOpdByKodeMaster godoc
// @Summary Get mapping opd by kode master
// @Description Get mapping opd by kode master
// @Tags mapping-opd
// @Produce json
// @Param kode_master query string true "Kode OPD master"
// @Success 200 {object} web.Response[opdWeb.MappingOpdResponse]
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /mapping-opd/find-by-kode-master [get]
func (app *Application) FindMappingOpdByKodeMaster(
	w http.ResponseWriter,
	r *http.Request,
) {
	q := r.URL.Query()

	kodeMaster := q.Get("kode_master")

	mapping, err := app.MappingOpd.FindByKodeMaster(
		r.Context(),
		kodeMaster,
	)
	if err != nil {
		if errors.Is(err, mappingOpdService.ErrMappingOpdNotFound) {
			app.NotFoundResponse(w, r)
			return
		}

		app.ServerErrorResponse(w, r, err)
		return
	}

	data := web.Response[*opdWeb.MappingOpdResponse]{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "DATA MAPPING OPD",
		Data:    mapping,
	}

	if err := app.WriteJSON(w, http.StatusOK, data, nil); err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

// FindMappingOpdByID godoc
// @Summary Get mapping opd by ID
// @Description Get mapping opd by ID
// @Tags mapping-opd
// @Produce json
// @Param id path int true "Mapping OPD ID"
// @Success 200 {object} web.Response[opdWeb.MappingOpdResponse]
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /mapping-opd/find/{id} [get]
func (app *Application) FindMappingOpdByID(
	w http.ResponseWriter,
	r *http.Request,
	params httprouter.Params,
) {
	id, err := readIDParam(params)
	if err != nil {
		app.BadRequestResponse(w, r, err.Error())
		return
	}

	mapping, err := app.MappingOpd.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, mappingOpdService.ErrMappingOpdNotFound) {
			app.NotFoundResponse(w, r)
			return
		}

		app.ServerErrorResponse(w, r, err)
		return
	}

	data := web.Response[*opdWeb.MappingOpdResponse]{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "DATA MAPPING OPD",
		Data:    mapping,
	}

	if err := app.WriteJSON(w, http.StatusOK, data, nil); err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

// CreateMappingOpd godoc
// @Summary Create mapping opd
// @Description Create mapping opd
// @Tags mapping-opd
// @Accept json
// @Produce json
// @Param request body opdWeb.CreateMappingOpdRequest true "Mapping OPD"
// @Success 201 {object} web.Response[opdWeb.MappingOpdResponse]
// @Failure 400 {object} web.ErrorResponse
// @Failure 409 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /mapping-opd/create [post]
func (app *Application) CreateMappingOpd(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req opdWeb.CreateMappingOpdRequest

	if err := app.ReadJSON(w, r, &req); err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	mapping := domain.MappingOpd{
		Sumber:     req.Sumber,
		KodeSumber: req.KodeSumber,
		KodeMaster: req.KodeMaster,
	}

	result, err := app.MappingOpd.Create(r.Context(), &mapping)
	if err != nil {
		switch {
		case errors.Is(err, mappingOpdService.ErrMappingOpdInvalid):
			app.BadRequestResponse(w, r, err.Error())

		case errors.Is(err, mappingOpdService.ErrMappingOpdExists):
			app.ConflictResponse(w, r, err.Error())

		default:
			app.ServerErrorResponse(w, r, err)
		}

		return
	}

	data := web.Response[*opdWeb.MappingOpdResponse]{
		Code:    http.StatusCreated,
		Status:  http.StatusText(http.StatusCreated),
		Message: "MAPPING OPD BERHASIL DIBUAT",
		Data:    result,
	}

	if err := app.WriteJSON(w, http.StatusCreated, data, nil); err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

// UpdateMappingOpd godoc
// @Summary Update mapping opd
// @Description Update mapping opd
// @Tags mapping-opd
// @Accept json
// @Produce json
// @Param id path int true "Mapping OPD ID"
// @Param request body opdWeb.UpdateMappingOpdRequest true "Mapping OPD"
// @Success 200 {object} web.Response[opdWeb.MappingOpdResponse]
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 409 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /mapping-opd/update/{id} [put]
func (app *Application) UpdateMappingOpd(
	w http.ResponseWriter,
	r *http.Request,
	params httprouter.Params,
) {
	id, err := readIDParam(params)
	if err != nil {
		app.BadRequestResponse(w, r, err.Error())
		return
	}

	var req opdWeb.UpdateMappingOpdRequest

	if err := app.ReadJSON(w, r, &req); err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	mapping := domain.MappingOpd{
		ID:         id,
		Sumber:     req.Sumber,
		KodeSumber: req.KodeSumber,
		KodeMaster: req.KodeMaster,
	}

	result, err := app.MappingOpd.Update(r.Context(), &mapping)
	if err != nil {
		switch {
		case errors.Is(err, mappingOpdService.ErrMappingOpdInvalid):
			app.BadRequestResponse(w, r, err.Error())

		case errors.Is(err, mappingOpdService.ErrMappingOpdNotFound):
			app.NotFoundResponse(w, r)

		case errors.Is(err, mappingOpdService.ErrMappingOpdExists):
			app.ConflictResponse(w, r, err.Error())

		default:
			app.ServerErrorResponse(w, r, err)
		}

		return
	}

	data := web.Response[*opdWeb.MappingOpdResponse]{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "MAPPING OPD BERHASIL DIUBAH",
		Data:    result,
	}

	if err := app.WriteJSON(w, http.StatusOK, data, nil); err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

// DeleteMappingOpd godoc
// @Summary Delete mapping opd
// @Description Delete mapping opd
// @Tags mapping-opd
// @Produce json
// @Param id path int true "Mapping OPD ID"
// @Success 200 {object} web.Response[any]
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /mapping-opd/delete/{id} [delete]
func (app *Application) DeleteMappingOpd(
	w http.ResponseWriter,
	r *http.Request,
	params httprouter.Params,
) {
	id, err := readIDParam(params)
	if err != nil {
		app.BadRequestResponse(w, r, err.Error())
		return
	}

	if err := app.MappingOpd.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, mappingOpdService.ErrMappingOpdInvalid):
			app.BadRequestResponse(w, r, err.Error())

		case errors.Is(err, mappingOpdService.ErrMappingOpdNotFound):
			app.NotFoundResponse(w, r)

		default:
			app.ServerErrorResponse(w, r, err)
		}

		return
	}

	data := web.Response[any]{
		Code:    http.StatusOK,
		Status:  http.StatusText(http.StatusOK),
		Message: "MAPPING OPD BERHASIL DIHAPUS",
		Data:    nil,
	}

	if err := app.WriteJSON(w, http.StatusOK, data, nil); err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}
}

func readIDParam(params httprouter.Params) (int64, error) {
	idString := params.ByName("id")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("id tidak valid")
	}

	return id, nil
}
