package handlers

import (
	"net/http"

	"github.com/shportix/blob-svc/internal/service/helpers"

	"github.com/shportix/blob-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeleteBlob(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewDeleteBlobRequest(r)
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
	allowed, _, err := helpers.CheckPermission(w, r, blob.OwnerAddress)
	if !allowed {
		return
	}
	err = helpers.BLobsQ(r).Delete(request.BlobID)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to delete blob from DB")
		ape.Render(w, problems.InternalError())
		return
	}

	w.WriteHeader(200)
}
