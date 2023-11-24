package handlers

import (
	"github.com/shportix/blob-svc/internal/data"
	"github.com/shportix/blob-svc/internal/service/helpers"
	"github.com/shportix/blob-svc/internal/service/requests"
	"github.com/shportix/blob-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/signcontrol"
	"net/http"
)

func GetBlobsList(w http.ResponseWriter, r *http.Request) {
	signer, err := signcontrol.CheckSignature(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(errors.Wrap(err, "Signature issues"))...)
		return
	}
	err = errors.New("Wrong signer")
	if signer != "GBA4EX43M25UPV4WIE6RRMQOFTWXZZRIPFAI5VPY6Z2ZVVXVWZ6NEOOB" {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
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
