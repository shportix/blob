package handlers

import (
	"gitlab.com/tokend/go/signcontrol"
	"net/http"

	"github.com/shportix/blob-svc/internal/service/helpers"

	"github.com/shportix/blob-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func DeleteBlob(w http.ResponseWriter, r *http.Request) {
	signer, err := signcontrol.CheckSignature(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(errors.Wrap(err, "Signature issues"))...)
		return
	}
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
	err = errors.New("Wrong signer")
	if signer != blob.OwnerAddress {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
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
