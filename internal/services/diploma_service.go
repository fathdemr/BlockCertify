package services

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/helper"
	"BlockCertify/internal/models"
	"BlockCertify/internal/repositories"
	apperrors "BlockCertify/pkg/errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type DiplomaService interface {
	Upload(filePath, fileHash string, metadata dto.DiplomaMetadataRequest) (*dto.UploadResponse, error)
	Verify(diplomaHash string) (*dto.VerifyResponse, error)
}

type diplomaService struct {
	arweave    ArweaveService
	blockchain BlockchainService
	repo       repositories.DiplomaRepository
}

func NewDiplomaService(arweave ArweaveService, blockchain BlockchainService, repo repositories.DiplomaRepository) DiplomaService {
	return &diplomaService{
		repo:       repo,
		arweave:    arweave,
		blockchain: blockchain,
	}
}

func (s *diplomaService) Upload(filePath, fileHash string, reqMeta dto.DiplomaMetadataRequest) (*dto.UploadResponse, error) {

	//Check before upload the arweave if diploma already exists
	log.Printf("Checking if diploma already exists: %s", fileHash)
	exists, existingArweaveTxID, err := s.blockchain.VerifyDiploma(fileHash)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrBlockchainFailed, "Failed to check if diploma exists", err)
	}

	if exists {
		return nil, apperrors.New(apperrors.ErrDiplomaExists, fmt.Sprintf("Diploma already registered. Arweave TxID: %s", existingArweaveTxID), err)
	}

	log.Println("Uploading to Arweave...")
	arweaveTxID, err := s.arweave.Upload(filePath, fileHash)
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

	log.Println("Storing on Polygon...")
	result, err := s.blockchain.StoreDiploma(fileHash, arweaveTxID)
	if err != nil {
		return nil, err
	}

	tx := s.repo.CreateTransaction()

	diplomaID := uuid.New()

	owner := fmt.Sprintf("%s %s", reqMeta.FirstName, reqMeta.LastName)

	diploma := models.Diploma{
		ID:          diplomaID,
		PublicID:    helper.GenerateDiplomaPublicIDFromUUID(diplomaID),
		Hash:        fileHash,
		ArweaveTxID: arweaveTxID,
		Owner:       owner,
		Timestamp:   time.Time{},
	}

	if err := tx.Create(&diploma).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	diplomaMetadata := models.DiplomaMetaData{
		ID:             uuid.New(),
		DiplomaID:      diploma.ID,
		FirstName:      reqMeta.FirstName,
		LastName:       reqMeta.LastName,
		Email:          reqMeta.Email,
		University:     reqMeta.University,
		Faculty:        reqMeta.Faculty,
		Department:     reqMeta.Department,
		GraduationYear: reqMeta.GraduationYear,
		StudentNumber:  reqMeta.StudentNumber,
		Nationality:    reqMeta.Nationality,
	}

	if err := tx.Create(&diplomaMetadata).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &dto.UploadResponse{
		Success:       true,
		DiplomaHash:   fileHash,
		ArweaveTxID:   arweaveTxID,
		ArweaveURL:    fmt.Sprintf("https://arweave.net/%s", arweaveTxID),
		PolygonTxHash: result.TransactionHash,
		BlockNumber:   result.BlockNumber,
	}, nil
}

func (s *diplomaService) Verify(diplomaHash string) (*dto.VerifyResponse, error) {
	log.Printf("%s - Verifying diploma on Polygon...", diplomaHash)

	exists, arweaveTxID, err := s.blockchain.VerifyDiploma(diplomaHash)
	if err != nil {
		return nil, err
	}

	response := &dto.VerifyResponse{
		Verified:    exists,
		DiplomaHash: diplomaHash,
	}

	if exists {
		response.ArweaveTxID = arweaveTxID
		response.ArweaveURL = fmt.Sprintf("https://arweave.net/%s", arweaveTxID)
	}

	return response, nil
}
