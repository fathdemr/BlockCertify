package dto

type VerifyResponse struct {
	Verified      bool   `json:"verified"`
	DiplomaHash   string `json:"diplomaHash"`
	ArweaveTxID   string `json:"arweaveTxID,omitempty"`
	ArweaveURL    string `json:"arweaveUrl,omitempty"`
	StudentName   string `json:"studentName,omitempty"`
	University    string `json:"university,omitempty"`
	Degree        string `json:"degree,omitempty"` // Represented by Department/Faculty in models
	IssueDate     string `json:"issueDate,omitempty"`
	PolygonTxHash string `json:"polygonTxHash,omitempty"`
	DiplomaID     string `json:"diplomaID"`
}
