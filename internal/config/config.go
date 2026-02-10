package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	Arweave    ArweaveConfig
	Blockchain BlockChainConfig
	JWTConfig  JWTConfig
	Db         DatabaseConfig
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

type JWTConfig struct {
	JWTExpireHours time.Duration
	JWTSecret      string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	chainIDStr := getEnvOrDefault("POLYGON_CHAIN_ID", "80002")
	chainID, err := strconv.Atoi(chainIDStr)
	if err != nil {
		return nil, fmt.Errorf("could not parse POLYGON_CHAIN_ID from env var: %w", err)
	}

	jwtExpStr := os.Getenv("JWT_EXP_HOURS")
	jwtExp, err := strconv.Atoi(jwtExpStr)
	if err != nil {
		return nil, fmt.Errorf("could not parse JWT_EXP_HOURS from env var: %w", err)
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
		JWTConfig: JWTConfig{
			JWTExpireHours: time.Duration(jwtExp),
			JWTSecret:      os.Getenv("JWT_SECRET_KEY"),
		},
		Db: DatabaseConfig{
			Host:     os.Getenv("APP_DB_HOST"),
			Port:     os.Getenv("APP_DB_PORT"),
			User:     os.Getenv("APP_DB_USERNAME"),
			Password: os.Getenv("APP_DB_PASSWORD"),
			Name:     os.Getenv("APP_DB_NAME"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("could not validate config: %w", err)
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
	if c.Db.Host == "" {
		return fmt.Errorf("APP_DB_HOST is required")
	}
	if c.Db.Port == "" {
		return fmt.Errorf("APP_DB_PORT is required")
	}
	if c.Db.User == "" {
		return fmt.Errorf("APP_DB_USERNAME is required")
	}
	if c.Db.Password == "" {
		return fmt.Errorf("APP_DB_PASSWORD is required")
	}
	if c.Db.Name == "" {
		return fmt.Errorf("APP_DB_NAME is required")
	}
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
