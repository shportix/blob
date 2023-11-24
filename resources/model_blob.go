package resources

import "encoding/json"

type Blob struct {
	Id           int64            `json:"id"`
	Data         json.RawMessage  `json:"data"`
	Relationship BlobRelationship `json:"relationship"`
}

type BlobRelationship struct {
	OwnerAddress string `json:"owner_id"`
}

type BlobResponse struct {
	Data Blob `json:"data"`
}

type BlobListResponse struct {
	Data []Blob `json:"data"`
}
