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
	admin  = "GBA4EX43M25UPV4WIE6RRMQOFTWXZZRIPFAI5VPY6Z2ZVVXVWZ6NEOOB"
	create = "create"
)

func CheckPermission(w http.ResponseWriter, r *http.Request, owner string) (bool, string, error) {
	signature := r.Header.Get("Signature")
	date, err := time.Parse(http.TimeFormat, r.Header.Get("date"))
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return false, "", err
	}
	realRequestTarget := r.Header.Get("Real-Request-Target")

	blobRequestsQ := BLobRequestsQ(r)
	blobRequests, err := blobRequestsQ.FilterBySign(signature).Get()
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return false, "", err
	}
	if blobRequests != nil && blobRequests.Date.Equal(date) && blobRequests.RealRequestTarget == realRequestTarget {
		ape.RenderErr(w, problems.Forbidden())
		return false, "", nil
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
		return false, "", err
	}

	err = blobRequestsQ.DeleteOld()
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return false, "", err
	}

	signer, err := signcontrol.CheckSignature(r)
	if err != nil {
		ape.RenderErr(w, problems.Forbidden())
		return false, signer, err
	}
	err = errors.New("Wrong signer")
	if owner == create || owner == signer || owner == admin {
		return true, signer, nil
	}
	ape.RenderErr(w, problems.Forbidden())
	return false, signer, err
}
