package domain

type SyncType string

const (
	SyncTypeOpd     SyncType = "OPD"
	SyncTypeJabatan SyncType = "JABATAN"
	SyncTypePegawai SyncType = "PEGAWAI"
)

func (j SyncType) IsValid() bool {

	switch j {
	case
		SyncTypeOpd:
		return true
	case
		SyncTypeJabatan:
		return true
	case
		SyncTypePegawai:
		return true
	}
	return false

}
