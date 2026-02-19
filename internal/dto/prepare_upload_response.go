package dto

type PrepareUploadResponse struct {
	DiplomaHash string `json:"diplomaHash"`
	ArweaveTxID string `json:"arweaveTxID"`
	ArweaveURL  string `json:"arweaveUrl"`
}
