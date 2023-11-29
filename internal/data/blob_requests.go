package data

import (
	"time"
)

type BlobRequestsQ interface {
	New() BlobRequestsQ

	Get() (*BlobRequest, error)

	Transaction(fn func(q BlobRequestsQ) error) error

	Insert(data BlobRequest) (BlobRequest, error)

	FilterBySign(sign ...string) BlobRequestsQ
	DeleteOld() error
}

type BlobRequest struct {
	ID                int64     `db:"id" structs:"-"`
	Sign              string    `db:"sign" structs:"sign"`
	RealRequestTarget string    `db:"real_request_target" structs:"real_request_target"`
	Date              time.Time `db:"date" structs:"date"`
}
