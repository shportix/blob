package data

import (
	"encoding/json"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type BlobsQ interface {
	New() BlobsQ

	Get() (*Blob, error)

	Delete(blobId int64) error
	Select() ([]Blob, error)

	Transaction(fn func(q BlobsQ) error) error

	Insert(data Blob) (Blob, error)

	Page(pageParams pgdb.OffsetPageParams) BlobsQ

	FilterByID(id ...int64) BlobsQ
}

type Blob struct {
	ID           int64           `db:"id" structs:"-"`
	Data         json.RawMessage `db:"data" structs:"data"`
	OwnerAddress string          `db:"owner_id" structs:"owner_id"`
}
