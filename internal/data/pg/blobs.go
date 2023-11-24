package pg

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/shportix/blob-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const blobsTableName = "blobs"

func NewBlobsQ(db *pgdb.DB) data.BlobsQ {
	return &blobsQ{
		db:  db.Clone(),
		sql: sq.Select("n.*").From(fmt.Sprintf("%s as n", blobsTableName)),
	}
}

type blobsQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

func (q *blobsQ) New() data.BlobsQ {
	return NewBlobsQ(q.db)
}

func (q *blobsQ) Get() (*data.Blob, error) {
	var result data.Blob
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *blobsQ) Delete(blobId int64) error {
	stmt := sq.Delete(blobsTableName).Where("id = ?", blobId)
	err := q.db.Exec(stmt)
	return err
}

func (q *blobsQ) Select() ([]data.Blob, error) {
	var result []data.Blob
	err := q.db.Select(&result, q.sql)
	return result, err
}

func (q *blobsQ) Transaction(fn func(q data.BlobsQ) error) error {
	return q.db.Transaction(func() error {
		return fn(q)
	})
}

func (q *blobsQ) Insert(value data.Blob) (data.Blob, error) {
	clauses := structs.Map(value)

	var result data.Blob
	stmt := sq.Insert(blobsTableName).SetMap(clauses).Suffix("returning *")
	err := q.db.Get(&result, stmt)

	return result, err
}

func (q *blobsQ) Page(pageParams pgdb.OffsetPageParams) data.BlobsQ {
	q.sql = pageParams.ApplyTo(q.sql, "id")
	return q
}

func (q *blobsQ) FilterByID(ids ...int64) data.BlobsQ {
	q.sql = q.sql.Where(sq.Eq{"n.id": ids})
	return q
}
