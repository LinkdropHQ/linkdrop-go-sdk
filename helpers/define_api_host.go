package helpers

import (
	"errors"
)

func DefineApiHost(apiUrl string, chainId int64) (string, error) {
	switch chainId {
	case 137:
		return apiUrl + "/polygon", nil
	case 8453:
		return apiUrl + "/base", nil
	case 42161:
		return apiUrl + "/arbitrum", nil
	case 10:
		return apiUrl + "/optimism", nil
	case 43114:
		return apiUrl + "/avalanche", nil
	}
	return "", errors.New("unsupported chain")
}
