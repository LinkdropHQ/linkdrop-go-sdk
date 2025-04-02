package linkdrop

type Option func(*SDKConfig, *ClientConfig)

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
		sdkc.environment = "production"
	}
}
