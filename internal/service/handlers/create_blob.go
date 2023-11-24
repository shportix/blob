package handlers

import (
	"github.com/shportix/blob-svc/internal/service/helpers"
	"gitlab.com/tokend/go/signcontrol"
	"net/http"

	"github.com/shportix/blob-svc/resources"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/shportix/blob-svc/internal/data"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/shportix/blob-svc/internal/service/requests"
)

func CreateBlob(w http.ResponseWriter, r *http.Request) {
	signer, err := signcontrol.CheckSignature(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(errors.Wrap(err, "Signature issues"))...)
		return
	}
	request, err := requests.NewCreateBlobRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	request.Data.Relationship.OwnerAddress = signer
	//err = errors.New("Wrong signer")
	//if request.Data.Relationship.OwnerAddress != signer {
	//	helpers.Log(r).WithError(err).Info("wrong request")
	//	ape.RenderErr(w, problems.BadRequest(err)...)
	//	return
	//}
	var resultBlob data.Blob

	err = helpers.BLobsQ(r).Transaction(func(q data.BlobsQ) error {
		blob := data.Blob{
			ID:           request.Data.Id,
			Data:         request.Data.Data,
			OwnerAddress: request.Data.Relationship.OwnerAddress,
		}

		resultBlob, err = q.Insert(blob)
		if err != nil {
			return errors.Wrap(err, "failed to insert blob")
		}

		return nil
	})
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	blobModel := resources.Blob{
		Id:   resultBlob.ID,
		Data: resultBlob.Data,
		Relationship: resources.BlobRelationship{
			OwnerAddress: resultBlob.OwnerAddress,
		},
	}

	result := resources.BlobResponse{
		Data: blobModel,
	}
	ape.Render(w, result)
}
