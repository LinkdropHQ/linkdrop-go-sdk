package linkdrop

import (
	"github.com/LinkdropHQ/linkdrop-go-sdk/constants"
	"github.com/ethereum/go-ethereum/common"
)

type Option func(*SDKConfig, *ClientConfig)

func WithEscrowContractAddress(escrowContractAddress common.Address) Option {
	return func(sdkc *SDKConfig, cc *ClientConfig) {
		sdkc.escrowContractAddress = escrowContractAddress
	}
}

func WithEscrowNFTContractAddress(escrowNFTContractAddress common.Address) Option {
	return func(sdkc *SDKConfig, cc *ClientConfig) {
		sdkc.escrowNFTContractAddress = escrowNFTContractAddress
	}
}

func WithEnvironmentTag(tag string) Option {
	return func(sdkc *SDKConfig, cc *ClientConfig) {
		sdkc.environment = tag
	}
}

func WithMessageConfig(messageConfig MessageConfig) Option {
	return func(sdkc *SDKConfig, cc *ClientConfig) {
		sdkc.messageConfig = messageConfig
	}
}

func WithApiUrl(apiUrl string) Option {
	return func(sdkc *SDKConfig, cc *ClientConfig) {
		cc.apiURL = apiUrl
	}
}

// Presets

func WithDefaultMessageConfig() Option {
	return func(sdkc *SDKConfig, cc *ClientConfig) {
		sdkc.applyDefaultMessageConfig()
	}
}

func WithProductionDefaults() Option {
	return func(sdkc *SDKConfig, cc *ClientConfig) {
		sdkc.applyProductionDefaults()
		sdkc.environment = "production"
	}
}

func WithCoinbaseWalletProductionDefaults() Option {
	return func(sdkc *SDKConfig, cc *ClientConfig) {
		sdkc.applyProductionDefaults()
		sdkc.escrowContractAddress = constants.CbwEscrowContractAddress
		sdkc.escrowNFTContractAddress = constants.CbwEscrowNFTContractAddress
		sdkc.environment = "production-coinbase-wallet"
	}
}
