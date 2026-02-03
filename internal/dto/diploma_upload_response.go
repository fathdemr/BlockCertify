package dto

type UploadResponse struct {
	Success       bool   `json:"success"`
	DiplomaHash   string `json:"diplomaHash"`
	ArweaveTxID   string `json:"arweaveTxID"`
	ArweaveURL    string `json:"arweaveUrl"`
	PolygonTxHash string `json:"polygonTxHash"`
	BlockNumber   uint64 `json:"blockNumber"`
}
