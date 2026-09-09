package domain

import "time"

type MappingOpd struct {
	ID         int64
	Sumber     string
	KodeSumber string
	KodeMaster string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
