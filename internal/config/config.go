package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	Arweave    ArweaveConfig
	Blockchain BlockChainConfig
}

type ServerConfig struct {
	Port          string
	UploadDir     string
	MaxUploadSize int64
}

type ArweaveConfig struct {
	WalletKey string
	Host      string
	Port      int
	Protocol  string
}

type BlockChainConfig struct {
	RPCURL          string
	PrivateKey      string
	ContractAddress string
	ChainID         int
	MinBalance      string // in MATIC
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	chainIDStr := getEnvOrDefault("POLYGON_CHAIN_ID", "80002")
	chainID, err := strconv.Atoi(chainIDStr)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:          getEnvOrDefault("PORT", "8080"),
			UploadDir:     getEnvOrDefault("UPLOAD_DIR", "upload"),
			MaxUploadSize: 10 << 20, // 10 MB
		},
		Arweave: ArweaveConfig{
			WalletKey: os.Getenv("ARWEAVE_KEY"),
			Host:      "arweave.net",
			Port:      443,
			Protocol:  "https",
		},
		Blockchain: BlockChainConfig{
			RPCURL:          getEnvOrDefault("POLYGON_RPC_URL", "https://polygon-rpc.com"),
			PrivateKey:      os.Getenv("PRIVATE_KEY"),
			ContractAddress: os.Getenv("CONTRACT_ADDRESS"),
			ChainID:         chainID,
			MinBalance:      "0.03",
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Arweave.WalletKey == "" {
		return fmt.Errorf("ARWEAVE_KEY is required")
	}
	if c.Blockchain.PrivateKey == "" {
		return fmt.Errorf("PRIVATE_KEY is required")
	}
	if c.Blockchain.ContractAddress == "" {
		return fmt.Errorf("CONTRACT_ADDRESS is required")
	}
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
