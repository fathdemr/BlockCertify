package dto

import "time"

type HistoryResponse struct {
	DiplomaID  string    `json:"diplomaId"`
	UserName   string    `json:"userName"`
	Department string    `json:"department"`
	CreateDate time.Time `json:"createDate"`
	DiplomaPdf string    `json:"diplomaPdf"`
}
