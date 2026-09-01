package domain

type SyncStatus string

const (
	SyncStatusSuccess SyncStatus = "SUCCESS"
)

func (j SyncStatus) IsValid() bool {

	switch j {
	case
		SyncStatusSuccess:
		return true
	}
	return false

}
