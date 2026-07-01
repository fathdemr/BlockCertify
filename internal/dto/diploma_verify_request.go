package dto

type VerifyDiplomaRequest struct {
	DiplomaID string `json:"diplomaId"`
	TxHash    string `json:"tx_hash"`
}
