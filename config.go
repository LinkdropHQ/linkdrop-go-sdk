package linkdrop

type MessageConfig struct {
	MinEncryptionKeyLength uint16
	MaxEncryptionKeyLength uint16
	MaxTextLength          int64
}

// ClientConfig is a configuration of the API Client
type ClientConfig struct {
	apiKey string
	apiURL string
}

type SDKConfig struct {
	baseURL       string
	messageConfig MessageConfig
	environment   string
}

// applyDefaults be run by SDK before any other options
func (sdkc *SDKConfig) applyDefaults() {
	sdkc.applyDefaultMessageConfig()
	sdkc.environment = "development"
}

func (sdkc *SDKConfig) applyDefaultMessageConfig() {
	sdkc.messageConfig = MessageConfig{
		MinEncryptionKeyLength: 6,
		MaxEncryptionKeyLength: 43,
		MaxTextLength:          140,
	}
}
