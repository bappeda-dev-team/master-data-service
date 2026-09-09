package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "master-data-service/docs"
)

func (app *Application) Routes() http.Handler {
	router := httprouter.New()

	// 404
	router.NotFound = http.HandlerFunc(app.NotFoundResponse)
	// 405
	router.MethodNotAllowed = http.HandlerFunc(app.MethodNotAllowedResponse)

	router.Handler(http.MethodGet, "/swagger/*filepath", httpSwagger.WrapHandler)

	router.HandlerFunc(http.MethodGet, "/v3/api-docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.HealthcheckHandler)

	// opd
	router.HandlerFunc(http.MethodGet, "/opd/findall", app.FindAllOpd)
	router.HandlerFunc(http.MethodPost, "/opd/sync", app.SyncDataMasterOpd)
	// mapping opd
	router.HandlerFunc(http.MethodGet, "/mapping-opd/findall", app.FindAllMappingOpd)
	router.HandlerFunc(http.MethodGet, "/mapping-opd/find-by-kode-sumber", app.FindMappingOpdByKodeSumber)
	router.HandlerFunc(http.MethodPost, "/mapping-opd/create", app.CreateMappingOpd)
	router.GET("/mapping-opd/find/:id", app.FindMappingOpdByID)
	router.PUT("/mapping-opd/update/:id", app.UpdateMappingOpd)
	router.DELETE("/mapping-opd/delete/:id", app.DeleteMappingOpd)

	// jabatan
	router.HandlerFunc(http.MethodGet, "/jabatan/find", app.FindJabatanByKodeOpd)
	router.HandlerFunc(http.MethodPost, "/jabatan/sync", app.SyncDataMasterJabatan)

	// pegawai
	router.HandlerFunc(http.MethodGet, "/pegawai/find", app.FindPegawaiByKodeOpd)
	router.HandlerFunc(http.MethodPost, "/pegawai/sync", app.SyncDataMasterPegawai)

	return app.recoverPanic(router)
}
