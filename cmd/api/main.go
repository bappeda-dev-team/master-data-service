package main

import (
	"log/slog"
	"os"
)

// @title           Master Data Service
// @version         1.0
// @description     Source Data Master. Not updateable, sync and view only
// @termsOfService  http://swagger.io/terms/
// @contact.name    Kertaskerja Dev Team
// @schemes http https
// @BasePath  /
func main() {

	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			nil,
		),
	)

	cfg := loadConfig(logger)

	// database
	db, err := openDB(cfg.DB.Dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(
				"failed closing database",
				"error",
				err,
			)
		}

	}()

	// application
	app := buildApplication(
		cfg,
		logger,
		db,
	)

	err = serve(app, cfg, logger)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
