package handlers

import (
	"net/http"

	"github.com/shportix/blob-svc/internal/service/helpers"

	"github.com/shportix/blob-svc/resources"

	"github.com/shportix/blob-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetBlob(w http.ResponseWriter, r *http.Request) {

	request, err := requests.NewGetBlobRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	blob, err := helpers.BLobsQ(r).FilterByID(request.BlobID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get blob from DB")
		ape.Render(w, problems.InternalError())
		return
	}

	allowed, _ := helpers.CheckPermission(w, r, blob.OwnerAddress)
	if !allowed {
		return
	}

	if blob == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	result := resources.BlobResponse{
		Data: resources.Blob{
			Data: blob.Data,
			Id:   blob.ID,
			Relationship: resources.BlobRelationship{
				OwnerAddress: blob.OwnerAddress,
			},
		},
	}

	ape.Render(w, result)
}
