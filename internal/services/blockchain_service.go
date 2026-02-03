package services

import (
	"BlockCertify/internal/config"
	"BlockCertify/internal/dto"
	"BlockCertify/internal/repositories"
	apperrors "BlockCertify/pkg/errors"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

type BlockchainService interface {
	StoreDiploma(diplomaHash, arweaveTxID string) (*dto.BlockchainResult, error)
	VerifyDiploma(diplomaHash string) (bool, string, error)
}

type blockchainService struct {
	repo       *repositories.ContractRepository
	minBalance *big.Int
	privateKey string
}

func NewBlockChainService(cfg *config.Config, repo *repositories.ContractRepository) BlockchainService {
	minBalance, _ := new(big.Float).SetString(cfg.Blockchain.MinBalance)
	minBalanceWei := new(big.Int)
	minBalance.Mul(minBalance, big.NewFloat(1e18)).Int(minBalanceWei)

	return &blockchainService{
		repo:       repo,
		minBalance: minBalanceWei,
		privateKey: cfg.Blockchain.PrivateKey,
	}
}

func (s *blockchainService) StoreDiploma(diplomaHash, arweaveTxID string) (*dto.BlockchainResult, error) {
	// Check if diploma already exists
	exists, _, err := s.repo.VerifyDiploma(diplomaHash)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrBlockchainFailed, "Failed to check diploma existence", err)
	}
	if exists {
		return nil, apperrors.New(apperrors.ErrDiplomaExists, "Diploma already registered in blockchain", nil)
	}

	// Check balance
	if err := s.checkBalance(); err != nil {
		return nil, err
	}

	// Get fee data for logging
	gasPrice, gasTipCap, err := s.repo.GetFeeData()
	if err != nil {
		return nil, apperrors.New(apperrors.ErrBlockchainFailed, "Failed to get fee data", err)
	}

	log.Printf("Using EIP-1559 fees: maxFeePerGas = %s gwei, maxPriorityFeePerGas = %s gwei",
		formatGwei(gasPrice),
		formatGwei(gasTipCap),
	)

	// Store diploma
	receipt, err := s.repo.StoreDiploma(diplomaHash, arweaveTxID)
	if err != nil {
		if err.Error() == "insufficient funds" {
			return nil, apperrors.New(
				apperrors.ErrInsufficientBalance,
				"Insufficient MATIC balance. Please get more test tokens from https://faucet.polygon.technology/",
				err,
			)
		}
		return nil, apperrors.New(apperrors.ErrBlockchainFailed, "Failed to store diploma", err)
	}

	if receipt.Status != 1 {
		return nil, apperrors.New(
			apperrors.ErrBlockchainFailed,
			fmt.Sprintf("Transaction reverted. Check on https://amoy.polygonscan.com/tx/%s", receipt.TxHash.Hex()),
			nil,
		)
	}

	log.Printf("Transaction confirmed in block %d", receipt.BlockNumber.Uint64())

	return &dto.BlockchainResult{
		TransactionHash: receipt.TxHash.Hex(),
		BlockNumber:     receipt.BlockNumber.Uint64(),
	}, nil
}

func (s *blockchainService) VerifyDiploma(diplomaHash string) (bool, string, error) {

	exists, arweaveTxID, err := s.repo.VerifyDiploma(diplomaHash)
	if err != nil {
		return false, "", apperrors.New(apperrors.ErrVerificationFailed, "Failed to verify diploma", err)
	}
	return exists, arweaveTxID, nil
}

func (s *blockchainService) checkBalance() error {
	privateKey, err := crypto.HexToECDSA(s.privateKey)
	if err != nil {
		return apperrors.New(apperrors.ErrBlockchainFailed, "Invalid private key", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return apperrors.New(apperrors.ErrBlockchainFailed, "Failed to cast public key to ECDSA", err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	balance, err := s.repo.GetBalance(fromAddress)
	if err != nil {
		return apperrors.New(apperrors.ErrBlockchainFailed, "Failed to get balance", err)
	}

	balanceInMatic := new(big.Float).Quo(new(big.Float).SetInt(balance), new(big.Float).SetFloat64(1e18))
	log.Printf("Wallet balance: %s MATIC", balanceInMatic.String())

	if balance.Cmp(s.minBalance) < 0 {
		return apperrors.New(apperrors.ErrInsufficientBalance, fmt.Sprintf("Insufficient MATIC balance. Current: %s MATIC. Please get more test tokens from https://faucet.polygon.technology/", balanceInMatic.String()), nil)
	}
	return nil
}

func formatGwei(wei *big.Int) string {
	gwei := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e9))
	return gwei.String()
}
