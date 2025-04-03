package helpers

import (
	"errors"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"net/url"
	"strconv"
	"strings"
)

func GetClaimCodeFromDashboardLink(claimUrl string) string {
	if strings.Contains(claimUrl, "redeem") {
		parts := strings.Split(claimUrl, "/")
		claimCode := strings.Split(parts[len(parts)-1], "?")[0] // TODO length check
		return claimCode
	} else {
		parsedUrl, err := url.Parse(claimUrl)
		if err != nil {
			return ""
		}
		queryParams := parsedUrl.Query()
		return queryParams.Get("k")
	}
}

func GetChainIdFromDashboardLink(claimUrl string) (chainId types.ChainId, err error) {
	parsedUrl, err := url.Parse(claimUrl)
	if err != nil {
		return
	}
	queryParams := parsedUrl.Query()
	cid, err := strconv.ParseInt(queryParams.Get("c"), 10, 64)
	if err != nil {
		return types.ChainIdPolygon, nil
	}
	chainId = types.ChainId(cid)
	if !chainId.IsSupported() {
		return chainId, errors.New("unsupported chain")
	}
	return
}
