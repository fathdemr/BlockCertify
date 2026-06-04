package BlockchainService

import (
	"BlockCertify/internal/dto"
	apperrors "BlockCertify/internal/pkg/errors"
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockchainService struct {
	minBalance      *big.Int
	RPCURL          string
	PrivateKey      string
	ContractAddress string
	ChainID         int
	MinBalance      string // in MATIC
	client          *ethclient.Client
	contractAddress common.Address
	contractABI     abi.ABI
}

const contractABI = `[
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "uint256",
        "name": "diplomaId",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "string",
        "name": "diplomaHash",
        "type": "string"
      },
      {
        "indexed": false,
        "internalType": "string",
        "name": "arweaveTxId",
        "type": "string"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "owner",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "timestamp",
        "type": "uint256"
      }
    ],
    "name": "DiplomaStored",
    "type": "event"
  },
  {
    "inputs": [],
    "name": "diplomaCount",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "name": "diplomas",
    "outputs": [
      {
        "internalType": "string",
        "name": "diplomaHash",
        "type": "string"
      },
      {
        "internalType": "string",
        "name": "arweaveTxId",
        "type": "string"
      },
      {
        "internalType": "address",
        "name": "owner",
        "type": "address"
      },
      {
        "internalType": "uint256",
        "name": "timestamp",
        "type": "uint256"
      },
      {
        "internalType": "bool",
        "name": "exists",
        "type": "bool"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "uint256",
        "name": "_diplomaId",
        "type": "uint256"
      }
    ],
    "name": "getDiploma",
    "outputs": [
      {
        "internalType": "string",
        "name": "diplomaHash",
        "type": "string"
      },
      {
        "internalType": "string",
        "name": "arweaveTxId",
        "type": "string"
      },
      {
        "internalType": "address",
        "name": "owner",
        "type": "address"
      },
      {
        "internalType": "uint256",
        "name": "timestamp",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "",
        "type": "string"
      }
    ],
    "name": "hashExists",
    "outputs": [
      {
        "internalType": "bool",
        "name": "",
        "type": "bool"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "",
        "type": "string"
      }
    ],
    "name": "hashToId",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "_diplomaHash",
        "type": "string"
      },
      {
        "internalType": "string",
        "name": "_arweaveTxId",
        "type": "string"
      }
    ],
    "name": "storeDiploma",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "_diplomaHash",
        "type": "string"
      }
    ],
    "name": "verifyDiploma",
    "outputs": [
      {
        "internalType": "bool",
        "name": "exists",
        "type": "bool"
      },
      {
        "internalType": "string",
        "name": "arweaveTxId",
        "type": "string"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  }
]`

func New(RPCURL, PrivateKey, ContractAddress, MinBalance string, ChainID int) *BlockchainService {
	minBalance, _ := new(big.Float).SetString(MinBalance)
	minBalanceWei := new(big.Int)
	minBalance.Mul(minBalance, big.NewFloat(1e18)).Int(minBalanceWei)
	client, err := ethclient.Dial(RPCURL)
	if err != nil {
		fmt.Println("failed to connect to blockchain: %w", err)
		return nil
	}

	parsedABI, err := abi.JSON(bytes.NewReader([]byte(contractABI)))
	if err != nil {
		fmt.Println("failed to parse contract ABI: %w", err)
		return nil
	}

	return &BlockchainService{
		RPCURL:          RPCURL,
		PrivateKey:      PrivateKey,
		ContractAddress: ContractAddress,
		ChainID:         ChainID,
		minBalance:      minBalanceWei,
		client:          client,
		contractAddress: common.HexToAddress(ContractAddress),
		contractABI:     parsedABI,
	}
}

func (s *BlockchainService) StoreDiploma(diplomaHash, arweaveTxID string) (*dto.BlockchainResult, error) {
	// Check if diploma already exists
	exists, _, err := s.VerifyDiplomaWithContract(diplomaHash)
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
	gasPrice, gasTipCap, err := s.GetFeeData()
	if err != nil {
		return nil, apperrors.New(apperrors.ErrBlockchainFailed, "Failed to get fee data", err)
	}

	log.Printf("Using EIP-1559 fees: maxFeePerGas = %s gwei, maxPriorityFeePerGas = %s gwei",
		formatGwei(gasPrice),
		formatGwei(gasTipCap),
	)

	// Store diploma
	receipt, err := s.StoreDiplomaWithContract(diplomaHash, arweaveTxID)
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

	slog.Info("Transaction confirmed in block %d", receipt.BlockNumber.Uint64())

	return &dto.BlockchainResult{
		TransactionHash: receipt.TxHash.Hex(),
		BlockNumber:     receipt.BlockNumber.Uint64(),
	}, nil
}

func (s *BlockchainService) VerifyDiploma(diplomaHash string) (bool, string, error) {

	exists, arweaveTxID, err := s.VerifyDiplomaWithContract(diplomaHash)
	if err != nil {
		return false, "", apperrors.New(apperrors.ErrVerificationFailed, "Failed to verify diploma", err)
	}
	return exists, arweaveTxID, nil
}

func (s *BlockchainService) checkBalance() error {
	privateKey, err := crypto.HexToECDSA(s.PrivateKey)
	if err != nil {
		return apperrors.New(apperrors.ErrBlockchainFailed, "Invalid private key", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return apperrors.New(apperrors.ErrBlockchainFailed, "Failed to cast public key to ECDSA", err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	balance, err := s.GetBalance(fromAddress)
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

func (s *BlockchainService) GetBalance(address common.Address) (*big.Int, error) {
	ctx := context.Background()
	return s.client.BalanceAt(ctx, address, nil)
}

func (s *BlockchainService) GetFeeData() (*big.Int, *big.Int, error) {

	ctx := context.Background()

	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, err
	}

	gasTipCap, err := s.client.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, nil, err
	}

	return gasPrice, gasTipCap, nil
}

func (s *BlockchainService) VerifyDiplomaWithContract(diplomaHash string) (bool, string, error) {

	ctx := context.Background()

	data, err := s.contractABI.Pack("verifyDiploma", diplomaHash)
	if err != nil {
		return false, "", err
	}

	result, err := s.client.CallContract(ctx, ethereum.CallMsg{
		To:   &s.contractAddress,
		Data: data,
	}, nil)
	if err != nil {
		return false, "", err
	}

	unpacked, err := s.contractABI.Unpack("verifyDiploma", result)
	if err != nil {
		return false, "", err
	}

	if len(unpacked) < 2 {
		return false, "", fmt.Errorf("Unexpected return values")
	}

	exists, ok := unpacked[0].(bool)
	if !ok {
		return false, "", fmt.Errorf("Failed to parse exists value")
	}

	arweaveTxID, ok := unpacked[1].(string)
	if !ok {
		return false, "", fmt.Errorf("Failed to parse arweaveTxId value")
	}
	return exists, arweaveTxID, nil
}

func (s *BlockchainService) StoreDiplomaWithContract(diplomaHash, arweaveTxID string) (*types.Receipt, error) {

	ctx := context.Background()

	privateKey, err := crypto.HexToECDSA(s.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to cast public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := s.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, gasTipCap, err := s.GetFeeData()
	if err != nil {
		return nil, err
	}

	data, err := s.contractABI.Pack("storeDiploma", diplomaHash, arweaveTxID)
	if err != nil {
		return nil, err
	}

	gasLimit, err := s.client.EstimateGas(ctx, ethereum.CallMsg{
		From: fromAddress,
		To:   &s.contractAddress,
		Data: data,
	})
	if err != nil {
		return nil, fmt.Errorf("gas estimation failed: %w", err)
	}

	// transfer int64 chainID to big.Int for EIP-1559
	chainID := big.NewInt(int64(s.ChainID))

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasPrice,
		Gas:       gasLimit,
		To:        &s.contractAddress,
		Value:     big.NewInt(0),
		Data:      data,
	})

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		return nil, err
	}

	if err := s.client.SendTransaction(ctx, signedTx); err != nil {
		return nil, err
	}

	receipt, err := s.waitForReceipt(signedTx.Hash())
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

func (s *BlockchainService) waitForReceipt(txHash common.Hash) (*types.Receipt, error) {

	ctx := context.Background()

	for i := 0; i < 60; i++ {
		receipt, err := s.client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}
		// Wait and retry
		select {
		case <-time.After(time.Second):
			continue
		}
	}
	return nil, fmt.Errorf("transaction receipt not found after timeout")
}

func (s *BlockchainService) GetAddress() string {
	privateKey, err := crypto.HexToECDSA(s.PrivateKey)
	if err != nil {
		return ""
	}
	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return ""
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
}

func (s *BlockchainService) Close() {
	if s.client != nil {
		s.client.Close()
	}
}
