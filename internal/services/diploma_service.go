package services

import (
	"BlockCertify/internal/models"
	apperrors "BlockCertify/pkg/errors"
	"fmt"
	"log"
)

type DiplomaService interface {
	Upload(filePath, fileHash string) (*models.UploadResponse, error)
	Verify(diplomaHash string) (*models.VerifyResponse, error)
}

type diplomaService struct {
	arweave    ArweaveService
	blockchain BlockchainService
}

func NewDiplomaService(arweave ArweaveService, blockchain BlockchainService) DiplomaService {
	return &diplomaService{
		arweave:    arweave,
		blockchain: blockchain,
	}
}

func (s *diplomaService) Upload(filePath, fileHash string) (*models.UploadResponse, error) {

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

	return &models.UploadResponse{
		Success:       true,
		DiplomaHash:   fileHash,
		ArweaveTxID:   arweaveTxID,
		ArweaveURL:    fmt.Sprintf("https://arweave.net/%s", arweaveTxID),
		PolygonTxHash: result.TransactionHash,
		BlockNumber:   result.BlockNumber,
	}, nil
}

func (s *diplomaService) Verify(diplomaHash string) (*models.VerifyResponse, error) {
	log.Printf("%s - Verifying diploma on Polygon...", diplomaHash)

	exists, arweaveTxID, err := s.blockchain.VerifyDiploma(diplomaHash)
	if err != nil {
		return nil, err
	}

	response := &models.VerifyResponse{
		Verified:    exists,
		DiplomaHash: diplomaHash,
	}

	if exists {
		response.ArweaveTxID = arweaveTxID
		response.ArweaveURL = fmt.Sprintf("https://arweave.net/%s", arweaveTxID)
	}

	return response, nil
}
