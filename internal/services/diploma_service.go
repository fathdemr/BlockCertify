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

	"github.com/gofrs/uuid/v5"
)

type DiplomaService interface {
	Upload(filePath, fileHash string, metadata dto.DiplomaMetadataRequest) (*dto.UploadResponse, error)
	Verify(req dto.VerifyDiplomaRequest) (dto.VerifyResponse, error)
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

	diplomaID, err := uuid.NewV7()

	owner := fmt.Sprintf("%s %s", reqMeta.FirstName, reqMeta.LastName)
	polygonURL := fmt.Sprintf("https://amoy.polygonscan.com/tx/%s", result.TransactionHash)
	arweaveURL := fmt.Sprintf("https://arweave.net/%s", arweaveTxID)

	diploma := models.Diploma{
		ID:          diplomaID,
		PublicID:    helper.GenerateDiplomaPublicIDFromUUID(diplomaID),
		Hash:        fileHash,
		ArweaveTxID: arweaveTxID,
		ArweaveURL:  arweaveURL,
		PolygonTxID: result.TransactionHash,
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
		ArweaveURL:    arweaveURL,
		PolygonTxHash: result.TransactionHash,
		BlockNumber:   result.BlockNumber,
	}, nil
}

func (s *diplomaService) ResolveHashFromReference(polygonTx, arweaveTxID string) (string, error) {
	log.Println("Resolving diploma hash from references...")

	// 1️⃣ Try Polygon reference first (strongest proof)
	if polygonTx != "" {
		log.Printf("Resolving hash from polygon tx: %s", polygonTx)

		hash, err := s.repo.GetHashFromPolygonTxID(polygonTx)
		if err == nil && hash != "" {
			return hash, nil
		}

		log.Printf("Polygon tx not found or error: %v", err)
	}

	// 2️⃣ Try Arweave reference
	if arweaveTxID != "" {
		log.Printf("Resolving hash from arweave tx: %s", arweaveTxID)

		hash, err := s.repo.GetHashFromArweaveTxID(arweaveTxID)
		if err == nil && hash != "" {
			return hash, nil
		}

		log.Printf("Arweave tx not found or error: %v", err)
	}

	return "", fmt.Errorf("could not resolve diploma hash from provided references")
}

func (s *diplomaService) Verify(req dto.VerifyDiplomaRequest) (dto.VerifyResponse, error) {

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
