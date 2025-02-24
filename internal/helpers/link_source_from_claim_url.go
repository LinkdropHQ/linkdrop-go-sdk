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
	src := types.CLSource(parsedUrl.Query().Get("src"))
	if src == types.CLSourceUndefined {
		return types.CLSourceP2P, nil
	}
	return src, nil
}
