package models

import "time"

type Diploma struct {
	Hash        string
	ArweaveTxID string
	Owner       string
	Timestamp   time.Time
}

type UploadRequest struct {
	File     []byte
	Filename string
}

type UploadResponse struct {
	Success       bool   `json:"success"`
	DiplomaHash   string `json:"diplomaHash"`
	ArweaveTxID   string `json:"arweaveTxID"`
	ArweaveURL    string `json:"arweaveUrl"`
	PolygonTxHash string `json:"polygonTxHash"`
	BlockNumber   uint64 `json:"blockNumber"`
}

type VerifyRequest struct {
	File     []byte
	Filename string
}

type VerifyResponse struct {
	Verified    bool   `json:"verified"`
	DiplomaHash string `json:"diplomaHash"`
	ArweaveTxID string `json:"arweaveTxID,omitempty"`
	ArweaveURL  string `json:"arweaveUrl,omitempty"`
}

type BlockchainResult struct {
	TransactionHash string
	BlockNumber     uint64
}
