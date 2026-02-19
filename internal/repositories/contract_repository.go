package repositories

import (
	"BlockCertify/internal/config"
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

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

type ContractRepository struct {
	client          *ethclient.Client
	contractAddress common.Address
	contractABI     abi.ABI
	privateKey      string
	chainID         *big.Int
}

func NewContractRepository(cfg *config.Config) (*ContractRepository, error) {
	client, err := ethclient.Dial(cfg.Blockchain.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to blockchain: %w", err)
	}

	//TODO : put a file reader for ABI and read from file instead of hardcoding
	parsedABI, err := abi.JSON(bytes.NewReader([]byte(contractABI)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	return &ContractRepository{
		client:          client,
		contractAddress: common.HexToAddress(cfg.Blockchain.ContractAddress),
		contractABI:     parsedABI,
		privateKey:      cfg.Blockchain.PrivateKey,
		chainID:         big.NewInt(int64(cfg.Blockchain.ChainID)),
	}, nil
}

func (r *ContractRepository) Close() {
	if r.client != nil {
		r.client.Close()
	}
}

func (r *ContractRepository) GetBalance(address common.Address) (*big.Int, error) {
	ctx := context.Background()
	return r.client.BalanceAt(ctx, address, nil)
}

func (r *ContractRepository) GetFeeData() (*big.Int, *big.Int, error) {

	ctx := context.Background()

	gasPrice, err := r.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, err
	}

	gasTipCap, err := r.client.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, nil, err
	}

	return gasPrice, gasTipCap, nil
}

func (r *ContractRepository) VerifyDiploma(diplomaHash string) (bool, string, error) {

	ctx := context.Background()

	data, err := r.contractABI.Pack("verifyDiploma", diplomaHash)
	if err != nil {
		return false, "", err
	}

	result, err := r.client.CallContract(ctx, ethereum.CallMsg{
		To:   &r.contractAddress,
		Data: data,
	}, nil)
	if err != nil {
		return false, "", err
	}

	unpacked, err := r.contractABI.Unpack("verifyDiploma", result)
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

func (r *ContractRepository) StoreDiploma(diplomaHash, arweaveTxID string) (*types.Receipt, error) {

	ctx := context.Background()

	privateKey, err := crypto.HexToECDSA(r.privateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to cast public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := r.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, gasTipCap, err := r.GetFeeData()
	if err != nil {
		return nil, err
	}

	data, err := r.contractABI.Pack("storeDiploma", diplomaHash, arweaveTxID)
	if err != nil {
		return nil, err
	}

	gasLimit, err := r.client.EstimateGas(ctx, ethereum.CallMsg{
		From: fromAddress,
		To:   &r.contractAddress,
		Data: data,
	})
	if err != nil {
		return nil, fmt.Errorf("gas estimation failed: %w", err)
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   r.chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasPrice,
		Gas:       gasLimit,
		To:        &r.contractAddress,
		Value:     big.NewInt(0),
		Data:      data,
	})

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(r.chainID), privateKey)
	if err != nil {
		return nil, err
	}

	if err := r.client.SendTransaction(ctx, signedTx); err != nil {
		return nil, err
	}

	receipt, err := r.waitForReceipt(signedTx.Hash())
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

func (r *ContractRepository) waitForReceipt(txHash common.Hash) (*types.Receipt, error) {

	ctx := context.Background()

	for i := 0; i < 60; i++ {
		receipt, err := r.client.TransactionReceipt(ctx, txHash)
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
