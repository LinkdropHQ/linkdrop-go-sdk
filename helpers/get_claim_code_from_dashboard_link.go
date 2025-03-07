package helpers

import (
	"net/url"
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
