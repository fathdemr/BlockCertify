package dto

// ConfirmUploadRequest is sent by the frontend after MetaMask has signed
// and submitted the Polygon transaction for a diploma.
type ConfirmUploadRequest struct {
	// Diploma identity (returned by /prepare)
	DiplomaHash string `json:"diplomaHash" binding:"required"`
	ArweaveTxID string `json:"arweaveTxID" binding:"required"`

	// Polygon tx data (provided by MetaMask / ethers.js)
	PolygonTxHash string `json:"polygonTxHash" binding:"required"`
	BlockNumber   uint64 `json:"blockNumber"`

	// Student metadata
	FirstName      string `json:"firstName" binding:"required"`
	LastName       string `json:"lastName" binding:"required"`
	Email          string `json:"email" binding:"required"`
	University     string `json:"university" binding:"required"`
	Faculty        string `json:"faculty"`
	Department     string `json:"department" binding:"required"`
	GraduationYear int    `json:"graduationYear" binding:"required"`
	StudentNumber  string `json:"studentNumber"`
	Nationality    string `json:"nationality"`
}
