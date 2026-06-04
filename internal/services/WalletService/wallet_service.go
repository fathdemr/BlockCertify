package WalletService

import (
	"errors"
	"log/slog"

	"github.com/everFinance/goar"
)

type WalletService struct {
	client *goar.Client
}

func New() *WalletService {

	client := goar.NewClient("https://arweave.net")

	return &WalletService{
		client: client,
	}
}

func (s *WalletService) ConnectWalletFromJSON(keyJSON []byte) (*goar.Wallet, error) {

	if len(keyJSON) == 0 {
		slog.Error("keyJSON is empty")
		return nil, errors.New("key is empty")
	}

	wallet, err := goar.NewWallet(keyJSON, "https://arweave.net")
	if err != nil {
		slog.Error("NewWalletFromJSON error:", err)
		return nil, err
	}

	slog.Info("wallet connected")
	return wallet, nil
}

func (s *WalletService) GetAddress(wallet *goar.Wallet) string {
	if wallet == nil {
		slog.Error("wallet is nil")
		return ""
	}
	return wallet.Signer.Address
}

func (s *WalletService) GetBalance(wallet *goar.Wallet) string {
	if wallet == nil {
		slog.Error("wallet is nil")
		return ""
	}

	balance, err := s.client.GetWalletBalance(wallet.Signer.Address)
	if err != nil {
		slog.Error("GetBalance error:", err)
		return ""
	}

	return balance.Text('f', 12)
}
