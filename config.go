package linkdrop

import (
	"github.com/LinkdropHQ/linkdrop-go-sdk/internal/constants"
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type MessageConfig struct {
	MinEncryptionKeyLength int64
	MaxEncryptionKeyLength int64
	MaxTextLength          int64
}

type Config struct {
	apiKey                   string
	baseURL                  string
	apiURL                   string
	dashboardURL             string
	timeout                  time.Duration
	retryCount               int64
	nativeTokenAddress       common.Address
	escrowContractAddress    common.Address
	escrowNFTContractAddress common.Address
	messageConfig            MessageConfig
	environment              string
}

func (c *Config) Environment() string {
	return c.environment
}

func (c *Config) applyDefaultMessageConfig() {
	c.messageConfig = MessageConfig{
		MinEncryptionKeyLength: 6,
		MaxEncryptionKeyLength: 43,
		MaxTextLength:          140,
	}
}

func (c *Config) applyDefaults() {
	c.apiURL = constants.ApiURL
	c.dashboardURL = constants.DevDashboardApiUrl
	c.timeout = 60 * time.Second
	c.retryCount = 5
	c.nativeTokenAddress = constants.NativeTokenAddress
	c.escrowContractAddress = constants.EscrowContractAddress
	c.escrowNFTContractAddress = constants.EscrowNFTContractAddress
	c.environment = "development"
	c.applyDefaultMessageConfig()
}

func (c *Config) applyProductionDefaults() {
	c.apiURL = constants.ApiURL
	c.dashboardURL = constants.DashboardApiUrl
	c.timeout = 10 * time.Second
	c.retryCount = 3
	c.nativeTokenAddress = constants.NativeTokenAddress
	c.escrowContractAddress = constants.EscrowContractAddress
	c.escrowNFTContractAddress = constants.EscrowNFTContractAddress
	c.applyDefaultMessageConfig()
}
