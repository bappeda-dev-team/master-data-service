package api

import (
	masterSyncDomain "master-data-service/internal/master_sync/domain"
)

func SyncMessage(status masterSyncDomain.SyncStatus) string {
	switch status {
	case masterSyncDomain.SyncStatusSuccess:
		return "Sinkronisasi OPD SIMPEG berhasil"

	case masterSyncDomain.SyncStatusPartial:
		return "Sinkronisasi OPD SIMPEG selesai dengan beberapa kegagalan"

	case masterSyncDomain.SyncStatusFailed:
		return "Sinkronisasi OPD SIMPEG gagal"

	default:
		return "Sinkronisasi OPD SIMPEG selesai"
	}
}
