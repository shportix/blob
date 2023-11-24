package requests

import (
	"encoding/json"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/shportix/blob-svc/resources"
)

type CreateBlobRequest struct {
	Data resources.Blob `json:"data"`
}

func NewCreateBlobRequest(r *http.Request) (CreateBlobRequest, error) {
	var request CreateBlobRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, nil
}

//func (r *CreateBlobRequest) validate() error {
//	return mergeErrors(validation.Errors{
//		"/data/data": validation.Validate(&r.Data.Data, validation.Required,
//			validation.Length(3, 10000)),
//	},
//	).Filter()
//}

//func mergeErrors(validationErrors ...validation.Errors) validation.Errors {
//	result := make(validation.Errors)
//	for _, errs := range validationErrors {
//		for key, err := range errs {
//			result[key] = err
//		}
//	}
//	return result
//}
