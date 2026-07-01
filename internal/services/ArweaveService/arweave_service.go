package ArweaveService

import (
	apperrors "BlockCertify/internal/pkg/errors"
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
)

type ArweaveService struct {
	client    *goar.Client
	wallet    *goar.Wallet
	WalletKey string
	Host      string
	Port      int
	Protocol  string
}

func NewArweaveService(WalletKey, Host, Protocol string, Port int) *ArweaveService {

	client := goar.NewClient("https://arweave.net")

	wallet, err := goar.NewWallet([]byte(WalletKey), "https://arweave.net")
	if err != nil {
		slog.Error("Failed to create Arweave wallet", "err", err)
	}

	return &ArweaveService{
		client:    client,
		wallet:    wallet,
		WalletKey: WalletKey,
		Host:      Host,
		Port:      Port,
		Protocol:  Protocol,
	}

}

func (s *ArweaveService) Upload(filePath, fileHash string) (string, error) {

	if s.wallet == nil {
		return "", apperrors.New(apperrors.ErrArweaveUploadFailed, "Arweave wallet not initialized — check walletKey in config", nil)
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", apperrors.New(apperrors.ErrArweaveUploadFailed, "Failed to read file", err)
	}

	// Check balance
	if err := s.checkBalanceForData(data); err != nil {
		return "", apperrors.New(apperrors.ErrArweaveUploadFailed, "Insufficient Arweave balance for transaction", err)
	}

	//Add tags
	tags := []types.Tag{
		{
			Name:  "Content-Type",
			Value: "application/pdf",
		},
		{
			Name:  "App-Name",
			Value: "DiplomaVerification",
		},
		{
			Name:  "File-Hash",
			Value: fileHash,
		},
		{
			Name:  "Timestamp",
			Value: strconv.FormatInt(time.Now().UnixMilli(), 10),
		},
	}

	// Create transaction
	tx, err := s.wallet.SendData(data, tags)
	if err != nil {
		return "", apperrors.New(apperrors.ErrArweaveUploadFailed, "Failed to upload to Arweave", err)
	}

	log.Printf("Arweave TxID: %s", tx.ID)
	return tx.ID, nil
}

func (s *ArweaveService) GetStatus() (address, balance string) {
	if s.wallet == nil {
		return "", ""
	}
	address = s.wallet.Signer.Address
	bal, err := s.client.GetWalletBalance(address)
	if err != nil {
		return address, "0"
	}
	return address, bal.Text('f', 6)
}

func (s *ArweaveService) checkBalanceForData(data []byte) error {
	if s.wallet == nil {
		return fmt.Errorf("arweave wallet not initialized")
	}
	balance, err := s.client.GetWalletBalance(s.wallet.Signer.Address)
	if err != nil {
		slog.Error("Failed to get wallet balance", "err", err)
	}

	//Estimate transaction cost based on data size(WINSTON)
	priceWinston, err := s.client.GetTransactionPrice(len(data), nil)
	if err != nil {
		slog.Error("Failed to get transaction price", "err", err)
	}

	//Convert winston -> AR
	priceAR := new(big.Float).SetInt64(priceWinston).Quo(new(big.Float).SetInt64(priceWinston), big.NewFloat(1e12))

	log.Printf("Arweave Balance: %s AR, estimated tx cost: %s AR", balance.Text('f', 12), priceAR.Text('f', 12))

	if balance.Cmp(priceAR) < 0 {
		slog.Error("Insufficient Arweave balance: required %s AR, available %s AR", priceAR.Text('f', 12), priceAR.Text('s', 12))
	}

	return nil
}
