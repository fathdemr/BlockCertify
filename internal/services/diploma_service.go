package services

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/helper"
	"BlockCertify/internal/models"
	apperrors "BlockCertify/internal/pkg/errors"
	"BlockCertify/internal/repositories"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

type DiplomaService interface {
	PrepareUpload(filePath, fileHash string, metadata dto.DiplomaMetadataRequest) (*dto.PrepareUploadResponse, error)
	ConfirmUpload(req dto.ConfirmUploadRequest) (*dto.UploadResponse, error)
	Verify(req dto.VerifyDiplomaRequest) (dto.VerifyResponse, error)
	GetArweaveUrlByDiplomaID(diplomaID string) string
	GetAllDiplomaFromDatabase() []dto.HistoryResponse
}

type diplomaService struct {
	Arweave    ArweaveService
	Blockchain BlockchainService
	repo       repositories.DiplomaRepository
}

func NewDiplomaService(arweave ArweaveService, blockchain BlockchainService, repo repositories.DiplomaRepository) DiplomaService {
	return &diplomaService{
		repo:       repo,
		Arweave:    arweave,
		Blockchain: blockchain,
	}
}

// PrepareUpload validates the diploma, uploads it to Arweave, and returns the
// hashes the frontend needs to sign the Polygon transaction via MetaMask.
// It does NOT touch the Polygon blockchain itself.
func (s *diplomaService) PrepareUpload(filePath, fileHash string, reqMeta dto.DiplomaMetadataRequest) (*dto.PrepareUploadResponse, error) {

	// Check on-chain if diploma already exists (read-only call, no signing)
	slog.Info("Checking if diploma exists on-chain", "hash", fileHash)
	exists, existingArweaveTxID, err := s.Blockchain.VerifyDiploma(fileHash)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrBlockchainFailed, "Failed to check if diploma exists", err)
	}

	if exists {
		return nil, apperrors.New(apperrors.ErrDiplomaExists, fmt.Sprintf("Diploma already registered. Arweave TxID: %s", existingArweaveTxID), nil)
	}

	slog.Info("Uploading to Arweave...")
	arweaveTxID, err := s.Arweave.Upload(filePath, fileHash)
	if err != nil {
		return nil, err
	}

	if arweaveTxID == "" {
		return nil, apperrors.New(
			apperrors.ErrArweaveUploadFailed,
			"Arweave upload failed: no transaction ID returned",
			nil,
		)
	}

	arweaveURL := fmt.Sprintf("https://arweave.net/%s", arweaveTxID)

	return &dto.PrepareUploadResponse{
		DiplomaHash: fileHash,
		ArweaveTxID: arweaveTxID,
		ArweaveURL:  arweaveURL,
	}, nil
}

// ConfirmUpload saves the diploma record to the database after the frontend
// has successfully submitted the Polygon transaction via MetaMask.
func (s *diplomaService) ConfirmUpload(req dto.ConfirmUploadRequest) (*dto.UploadResponse, error) {

	slog.Info("Confirming diploma upload", "polygonTxHash", req.PolygonTxHash)

	polygonURL := fmt.Sprintf("https://amoy.polygonscan.com/tx/%s", req.PolygonTxHash)
	arweaveURL := fmt.Sprintf("https://arweave.net/%s", req.ArweaveTxID)

	tx := s.repo.CreateTransaction()

	diplomaID, err := uuid.NewV7()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to generate diploma ID: %w", err)
	}

	owner := fmt.Sprintf("%s %s", req.FirstName, req.LastName)

	diploma := models.Diploma{
		ID:          diplomaID,
		PublicID:    helper.GenerateDiplomaPublicIDFromUUID(diplomaID),
		Hash:        req.DiplomaHash,
		ArweaveTxID: req.ArweaveTxID,
		ArweaveURL:  arweaveURL,
		PolygonTxID: req.PolygonTxHash,
		PolygonURL:  polygonURL,
		Owner:       owner,
		Timestamp:   time.Time{},
	}

	if err := tx.Create(&diploma).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	diplomaMetadata := models.DiplomaMetaData{
		ID:             uuid.Must(uuid.NewV7()),
		DiplomaID:      diploma.ID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		University:     req.University,
		Faculty:        req.Faculty,
		Department:     req.Department,
		GraduationYear: req.GraduationYear,
		StudentNumber:  req.StudentNumber,
		Nationality:    req.Nationality,
	}

	if err := tx.Create(&diplomaMetadata).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &dto.UploadResponse{
		Success:       true,
		DiplomaHash:   req.DiplomaHash,
		ArweaveTxID:   req.ArweaveTxID,
		ArweaveURL:    arweaveURL,
		PolygonTxHash: req.PolygonTxHash,
		BlockNumber:   req.BlockNumber,
	}, nil
}

func (s *diplomaService) Verify(req dto.VerifyDiplomaRequest) (dto.VerifyResponse, error) {

	slog.Info("Verifying diploma from diplomaID")

	var err error
	var response dto.VerifyResponse
	diplomaID := req.DiplomaID

	response.Verified = false

	// Fetch metadata from DB if it exists
	diploma, err := s.repo.GetByDiplomaID(diplomaID)

	if err == nil && diploma != nil {
		response.Verified = true
		response.StudentName = diploma.Owner
		if diploma.MetaData.ID != uuid.Nil {
			response.University = diploma.MetaData.University
			response.Degree = fmt.Sprintf("%s - %s", diploma.MetaData.Faculty, diploma.MetaData.Department)
			response.IssueDate = diploma.CreatedAt.Format("2006-01-02")
		}
		response.PolygonTxHash = diploma.PolygonTxID
		response.ArweaveTxID = diploma.ArweaveTxID
		response.ArweaveURL = diploma.ArweaveURL
		response.DiplomaHash = diploma.Hash
		response.DiplomaID = diplomaID
	}

	if !response.Verified {
		return dto.VerifyResponse{
			Verified: false,
		}, fmt.Errorf("diploma not found")
	}

	return response, nil
}

func (s *diplomaService) GetArweaveUrlByDiplomaID(diplomaID string) string {

	slog.Info("Fetching diploma by tx id")

	if strings.TrimSpace(diplomaID) == "" {
		slog.Error("Diploma ID is empty")
		return ""
	}

	diploma, err := s.repo.GetByDiplomaID(diplomaID)
	if err != nil {
		slog.Error("Failed to get diploma by tx id", "err", err)
		return ""
	}

	arweaveUrl := diploma.ArweaveURL
	if arweaveUrl == "" {
		slog.Error("Arweave URL is empty")
		return ""
	}

	return arweaveUrl
}

func (s *diplomaService) GetAllDiplomaFromDatabase() []dto.HistoryResponse {

	ch := s.repo.GetAllDiplomaFromDatabase()

	var response []dto.HistoryResponse

	for d := range ch {

		resp := dto.HistoryResponse{
			DiplomaID:  d.PublicID,
			UserName:   d.Owner,
			Department: d.MetaData.Department,
			CreateDate: d.CreatedAt,
			DiplomaPdf: d.ArweaveURL,
		}

		response = append(response, resp)
	}

	return response
}
