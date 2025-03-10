package linkdrop

import (
	"github.com/LinkdropHQ/linkdrop-go-sdk/constants"
	"github.com/ethereum/go-ethereum/common"
)

type Option func(*SDKConfig)

func WithEscrowContractAddress(escrowContractAddress common.Address) Option {
	return func(sdkc *SDKConfig) {
		sdkc.escrowContractAddress = escrowContractAddress
	}
}

func WithEscrowNFTContractAddress(escrowNFTContractAddress common.Address) Option {
	return func(sdkc *SDKConfig) {
		sdkc.escrowNFTContractAddress = escrowNFTContractAddress
	}
}

func WithEnvironmentTag(tag string) Option {
	return func(sdkc *SDKConfig) {
		sdkc.environment = tag
	}
}

func WithMessageConfig(messageConfig MessageConfig) Option {
	return func(sdkc *SDKConfig) {
		sdkc.messageConfig = messageConfig
	}
}

// Presets

func WithDefaultMessageConfig() Option {
	return func(sdkc *SDKConfig) {
		sdkc.applyDefaultMessageConfig()
	}
}

func WithProductionDefaults() Option {
	return func(sdkc *SDKConfig) {
		sdkc.applyProductionDefaults()
		sdkc.environment = "production"
	}
}

func WithCoinbaseWalletProductionDefaults() Option {
	return func(sdkc *SDKConfig) {
		sdkc.applyProductionDefaults()
		sdkc.escrowContractAddress = constants.CbwEscrowContractAddress
		sdkc.escrowNFTContractAddress = constants.CbwEscrowNFTContractAddress
		sdkc.environment = "production-coinbase-wallet"
	}
}
