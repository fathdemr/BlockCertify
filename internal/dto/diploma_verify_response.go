package dto

type VerifyResponse struct {
	Verified    bool   `json:"verified"`
	DiplomaHash string `json:"diplomaHash"`
	ArweaveTxID string `json:"arweaveTxID,omitempty"`
	ArweaveURL  string `json:"arweaveUrl,omitempty"`
}
