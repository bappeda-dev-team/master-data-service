package domain

type SyncStatus string

const (
	SyncStatusRunning SyncStatus = "RUNNING"
	SyncStatusSuccess SyncStatus = "SUCCESS"
	SyncStatusPartial SyncStatus = "PARTIAL"
	SyncStatusFailed  SyncStatus = "FAILED"
)

func (j SyncStatus) IsValid() bool {

	switch j {
	case
		SyncStatusRunning:
		return true
	case
		SyncStatusSuccess:
		return true
	case
		SyncStatusPartial:
		return true
	case
		SyncStatusFailed:
		return true
	}
	return false

}
