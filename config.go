package linkdrop

import (
	"github.com/LinkdropHQ/linkdrop-go-sdk/constants"
	"github.com/ethereum/go-ethereum/common"
)

type MessageConfig struct {
	MinEncryptionKeyLength int64
	MaxEncryptionKeyLength int64
	MaxTextLength          int64
}

// ClientConfig is a configuration of the API Client
type ClientConfig struct {
	apiKey string
	apiURL string
}

type SDKConfig struct {
	baseURL                  string
	escrowContractAddress    common.Address
	escrowNFTContractAddress common.Address
	messageConfig            MessageConfig
	environment              string
}

func (sdkc *SDKConfig) applyDefaults() {
	sdkc.escrowContractAddress = constants.EscrowContractAddress
	sdkc.escrowNFTContractAddress = constants.EscrowNFTContractAddress
	sdkc.environment = "development"
	sdkc.applyDefaultMessageConfig()
}

func (sdkc *SDKConfig) applyProductionDefaults() {
	sdkc.escrowContractAddress = constants.EscrowContractAddress
	sdkc.escrowNFTContractAddress = constants.EscrowNFTContractAddress
	sdkc.applyDefaultMessageConfig()
}

func (sdkc *SDKConfig) applyDefaultMessageConfig() {
	sdkc.messageConfig = MessageConfig{
		MinEncryptionKeyLength: 6,
		MaxEncryptionKeyLength: 43,
		MaxTextLength:          140,
	}
}
