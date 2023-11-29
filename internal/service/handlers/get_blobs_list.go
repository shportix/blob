package handlers

import (
	"github.com/shportix/blob-svc/internal/data"
	"github.com/shportix/blob-svc/internal/service/helpers"
	"github.com/shportix/blob-svc/internal/service/requests"
	"github.com/shportix/blob-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

const (
	admin = "GBA4EX43M25UPV4WIE6RRMQOFTWXZZRIPFAI5VPY6Z2ZVVXVWZ6NEOOB"
)

func GetBlobsList(w http.ResponseWriter, r *http.Request) {
	allowed, _, err := helpers.CheckPermission(w, r, admin)
	if !allowed {
		return
	}
	_, err = requests.NewGetBlobsListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	blobsQ := helpers.BLobsQ(r)
	blobs, err := blobsQ.Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get blobs")
		ape.Render(w, problems.InternalError())
		return
	}

	response := resources.BlobListResponse{
		Data: newBlobsList(blobs),
	}
	ape.Render(w, response)
}

func newBlobsList(blobs []data.Blob) []resources.Blob {
	result := make([]resources.Blob, len(blobs))
	for i, blob := range blobs {
		result[i] = resources.Blob{
			Id:   blob.ID,
			Data: blob.Data,
			Relationship: resources.BlobRelationship{
				OwnerAddress: blob.OwnerAddress,
			},
		}
	}
	return result
}
