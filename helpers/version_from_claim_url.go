package helpers

import (
	"errors"
	"net/url"
)

func VersionFromClaimUrl(claimUrl string) (string, error) {
	parsedParams, err := url.ParseQuery(claimUrl)
	if err != nil {
		return "", err
	}

	version := parsedParams.Get("v")
	if version == "" {
		return "", errors.New("version not provided")
	}

	return version, nil
}
