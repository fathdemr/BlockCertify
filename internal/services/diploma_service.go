package services

import (
	"BlockCertify/internal/dto"
	"BlockCertify/internal/helper"
	"BlockCertify/internal/models"
	"BlockCertify/internal/repositories"
	apperrors "BlockCertify/pkg/errors"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

type DiplomaService interface {
	Upload(filePath, fileHash string, metadata dto.DiplomaMetadataRequest) (*dto.UploadResponse, error)
	Verify(req dto.VerifyDiplomaRequest) (dto.VerifyResponse, error)
	GetArweaveUrlByDiplomaID(diplomaID string) string
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

func (s *diplomaService) Upload(filePath, fileHash string, reqMeta dto.DiplomaMetadataRequest) (*dto.UploadResponse, error) {

	//Check before upload the arweave if diploma already exists
	slog.Info("Checking if diploma exists", "hash", fileHash)
	exists, existingArweaveTxID, err := s.Blockchain.VerifyDiploma(fileHash)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrBlockchainFailed, "Failed to check if diploma exists", err)
	}

	if exists {
		return nil, apperrors.New(apperrors.ErrDiplomaExists, fmt.Sprintf("Diploma already registered. Arweave TxID: %s", existingArweaveTxID), err)
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

	log.Println("Storing on Polygon...")
	result, err := s.Blockchain.StoreDiploma(fileHash, arweaveTxID)
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
	slog.Info("Resolving hash from reference")

	// 1️⃣ Try Polygon reference first (strongest proof)
	if polygonTx != "" {
		slog.Info("Resolving hash from polygon tx", "tx", polygonTx)

		hash, err := s.repo.GetHashFromPolygonTxID(polygonTx)
		if err == nil && hash != "" {
			return hash, nil
		}

		slog.Info("Polygon tx not found or error", "err", err)
	}

	// 2️⃣ Try Arweave reference
	if arweaveTxID != "" {
		slog.Info("Resolving hash from arweave tx", "tx", arweaveTxID)

		hash, err := s.repo.GetHashFromArweaveTxID(arweaveTxID)
		if err == nil && hash != "" {
			return hash, nil
		}

		slog.Info("Arweave tx not found or error", "err", err)
	}

	err := fmt.Errorf("could not resolve diploma hash from provided references")

	slog.Error("could not resolve diploma hash from provided references", "polygonTx", polygonTx, "arweaveTxID", arweaveTxID, "err", err)

	return "", err

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
