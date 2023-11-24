package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"

	"gitlab.com/distributed_lab/urlval"
)

type DeleteBlobRequest struct {
	BlobID int64 `url:"-"`
}

func NewDeleteBlobRequest(r *http.Request) (DeleteBlobRequest, error) {
	request := DeleteBlobRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.BlobID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
