package pg

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/shportix/blob-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"time"
)

const blobRequestsTableName = "blobs_requests"

func NewBlobRequestsQ(db *pgdb.DB) data.BlobRequestsQ {
	return &blobRequestsQ{
		db:  db.Clone(),
		sql: sq.Select("n.*").From(fmt.Sprintf("%s as n", blobRequestsTableName)),
	}
}

type blobRequestsQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

func (q *blobRequestsQ) New() data.BlobRequestsQ {
	return NewBlobRequestsQ(q.db)
}

func (q *blobRequestsQ) Get() (*data.BlobRequest, error) {
	var result data.BlobRequest
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *blobRequestsQ) DeleteOld() error {
	stmt := sq.Delete(blobRequestsTableName).Where("date < ?", time.Now().UTC().Add(-1*time.Hour))
	err := q.db.Exec(stmt)
	return err
}

func (q *blobRequestsQ) Delete(blobRequestId int64) error {
	stmt := sq.Delete(blobRequestsTableName).Where("id = ?", blobRequestId)
	err := q.db.Exec(stmt)
	return err
}

func (q *blobRequestsQ) Select() ([]data.BlobRequest, error) {
	var result []data.BlobRequest
	err := q.db.Select(&result, q.sql)
	return result, err
}

func (q *blobRequestsQ) Transaction(fn func(q data.BlobRequestsQ) error) error {
	return q.db.Transaction(func() error {
		return fn(q)
	})
}

func (q *blobRequestsQ) Insert(value data.BlobRequest) (data.BlobRequest, error) {
	clauses := structs.Map(value)

	var result data.BlobRequest
	stmt := sq.Insert(blobRequestsTableName).SetMap(clauses).Suffix("returning *")
	err := q.db.Get(&result, stmt)

	return result, err
}

func (q *blobRequestsQ) Page(pageParams pgdb.OffsetPageParams) data.BlobRequestsQ {
	q.sql = pageParams.ApplyTo(q.sql, "id")
	return q
}

func (q *blobRequestsQ) FilterBySign(signs ...string) data.BlobRequestsQ {
	q.sql = q.sql.Where(sq.Eq{"n.sign": signs})
	return q
}
