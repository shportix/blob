package helpers

import (
	"github.com/shportix/blob-svc/internal/data"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/signcontrol"
	"net/http"
	"time"
)

const (
	admin                = "GBA4EX43M25UPV4WIE6RRMQOFTWXZZRIPFAI5VPY6Z2ZVVXVWZ6NEOOB"
	create               = "create"
	signatureHeader      = "Signature"
	realTimeTargetHeader = "Real-Request-Target"
)

func CheckPermission(w http.ResponseWriter, r *http.Request, owner string) (bool, string) {
	signature := r.Header.Get(signatureHeader)
	date, err := time.Parse(http.TimeFormat, r.Header.Get("date"))
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return false, ""
	}
	realRequestTarget := r.Header.Get(realTimeTargetHeader)

	blobRequestsQ := BLobRequestsQ(r)
	blobRequests, err := blobRequestsQ.FilterBySign(signature).Get()
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return false, ""
	}
	if blobRequests != nil && blobRequests.Date.Equal(date) && blobRequests.RealRequestTarget == realRequestTarget {
		ape.RenderErr(w, problems.Forbidden())
		return false, ""
	}

	err = BLobRequestsQ(r).Transaction(func(q data.BlobRequestsQ) error {
		blobRequest := data.BlobRequest{
			ID:                1,
			Sign:              signature,
			Date:              date,
			RealRequestTarget: realRequestTarget,
		}

		_, err = q.Insert(blobRequest)
		if err != nil {
			return errors.Wrap(err, "failed to insert blob")
		}

		return nil
	})
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return false, ""
	}

	err = blobRequestsQ.DeleteOld()
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return false, ""
	}

	signer, err := signcontrol.CheckSignature(r)
	if err != nil {
		ape.RenderErr(w, problems.Forbidden())
		return false, signer
	}
	err = errors.New("Wrong signer")
	if owner == create || signer == signer || signer == admin {
		return true, signer
	}
	Log(r).Info("hear")
	ape.RenderErr(w, problems.Forbidden())
	return false, signer
}
