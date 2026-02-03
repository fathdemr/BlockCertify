package dto

type DiplomaMetadataRequest struct {
	FirstName      string `json:"firstName" binding:"required"`
	LastName       string `json:"lastName" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	University     string `json:"university" binding:"required"`
	Faculty        string `json:"faculty" binding:"required"`
	Department     string `json:"department" binding:"required"`
	GraduationYear int    `json:"graduationYear" binding:"required"`
	StudentNumber  string `json:"studentNumber" binding:"required"`
	Nationality    string `json:"nationality" binding:"required"`
}
