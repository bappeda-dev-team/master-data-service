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

	return app.recoverPanic(router)
}
