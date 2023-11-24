package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type GetBlobsListRequest struct {
	pgdb.OffsetPageParams
	FilterId []string `filter:"id"`
}

func NewGetBlobsListRequest(r *http.Request) (GetBlobsListRequest, error) {
	request := GetBlobsListRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
