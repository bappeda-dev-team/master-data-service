package domain

type SyncType string

const (
	SyncTypeOpd SyncType = "OPD"
)

func (j SyncType) IsValid() bool {

	switch j {
	case
		SyncTypeOpd:
		return true
	}
	return false

}
