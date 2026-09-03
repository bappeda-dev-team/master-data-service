package api

import (
	"log/slog"
	jabatan "master-data-service/internal/jabatan/service"
	masterSync "master-data-service/internal/master_sync/service"
	opd "master-data-service/internal/opd/service"
	pegawai "master-data-service/internal/pegawai/service"
)

const Version = "1.0.0"

// port: which port of server should be listening to
// env: which env should used, prod, staging, and dev
type Config struct {
	Port int
	Env  string

	DB struct {
		Dsn string
	}

	Simpeg ServiceConfig
}

// ServiceConfig used for external services
type ServiceConfig struct {
	BaseURL string
	ApiPath string
}

// hold application dependencies for
// http handler, middleware and helpers
type Application struct {
	Config Config
	Logger *slog.Logger

	MasterSync *masterSync.MasterSyncService

	Opd *opd.DataMasterOpdService

	Jabatan *jabatan.MasterDataJabatanService

	Pegawai *pegawai.MasterDataPegawaiService
}
