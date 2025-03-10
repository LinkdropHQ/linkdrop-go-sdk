package helpers

import (
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"net/url"
)

func LinkSourceFromClaimUrl(claimUrl string) (types.LinkSource, error) {
	parsedUrl, err := url.Parse(claimUrl)
	if err != nil {
		return types.LinkSourceUndefined, err
	}
	src := types.LinkSource(parsedUrl.Query().Get("src"))
	if src == types.LinkSourceUndefined {
		return types.LinkSourceP2P, nil
	}
	return src, nil
}
