package domain

type SyncType string

const (
	SyncTypeOpd     SyncType = "OPD"
	SyncTypeJabatan SyncType = "JABATAN"
)

func (j SyncType) IsValid() bool {

	switch j {
	case
		SyncTypeOpd:
		return true
	case
		SyncTypeJabatan:
		return true
	}
	return false

}
