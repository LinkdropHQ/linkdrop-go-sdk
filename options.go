package linkdrop

import (
	"github.com/LinkdropHQ/linkdrop-go-sdk/internal/constants"
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type Option func(*Config)

func WithApiKey(key string) Option {
	return func(c *Config) {
		c.apiKey = key
	}
}

func WithApiURL(url string) Option {
	return func(c *Config) {
		c.apiURL = url
	}
}

func WithDashboardURL(url string) Option {
	return func(c *Config) {
		c.dashboardURL = url
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.timeout = timeout
	}
}

func WithRetryCount(retryCount int64) Option {
	return func(c *Config) {
		c.retryCount = retryCount
	}
}

func WithNativeTokenAddress(tokenAddress common.Address) Option {
	return func(c *Config) {
		c.nativeTokenAddress = tokenAddress
	}
}

func WithEscrowContractAddress(escrowContractAddress common.Address) Option {
	return func(c *Config) {
		c.escrowContractAddress = escrowContractAddress
	}
}

func WithEscrowNFTContractAddress(escrowNFTContractAddress common.Address) Option {
	return func(c *Config) {
		c.escrowNFTContractAddress = escrowNFTContractAddress
	}
}

func WithEnvironmentTag(tag string) Option {
	return func(c *Config) {
		c.environment = tag
	}
}

func WithMessageConfig(messageConfig MessageConfig) Option {
	return func(c *Config) {
		c.messageConfig = messageConfig
	}
}

// Presets

func WithDefaultMessageConfig() Option {
	return func(c *Config) {
		c.applyDefaultMessageConfig()
	}
}

func WithProductionDefaults() Option {
	return func(c *Config) {
		c.applyProductionDefaults()
		c.environment = "production"
	}
}

func WithCoinbaseWalletProductionDefaults() Option {
	return func(c *Config) {
		c.applyProductionDefaults()
		c.escrowContractAddress = constants.CbwEscrowContractAddress
		c.escrowNFTContractAddress = constants.CbwEscrowNFTContractAddress
		c.environment = "production-coinbase-wallet"
	}
}
