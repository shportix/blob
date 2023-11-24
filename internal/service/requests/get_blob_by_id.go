package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"

	"gitlab.com/distributed_lab/urlval"
)

type GetBlobRequest struct {
	BlobID int64 `url:"-"`
}

func NewGetBlobRequest(r *http.Request) (GetBlobRequest, error) {
	request := GetBlobRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.BlobID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
