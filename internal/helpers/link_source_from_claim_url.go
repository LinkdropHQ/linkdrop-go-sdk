package helpers

import (
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"net/url"
)

func LinkSourceFromClaimUrl(claimUrl string) (types.CLSource, error) {
	parsedUrl, err := url.Parse(claimUrl)
	if err != nil {
		return "unspecified", err
	}
	if parsedUrl.Query().Get("src") == "" {
		return types.CLSourceP2P, nil
	}
	return types.CLSourceD, nil
}
