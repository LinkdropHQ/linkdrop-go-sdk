package types

import (
	"errors"
	"github.com/LinkdropHQ/linkdrop-go-sdk/internal/constants"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type AuthorizationMethod string

const (
	AMApproveWithAuthorization    AuthorizationMethod = "ApproveWithAuthorization"
	AMReceiveWithAuthorization    AuthorizationMethod = "ReceiveWithAuthorization"
	AMReceiveWithAuthorizationEOA AuthorizationMethod = "ReceiveWithAuthorizationEOA"
)

func (am AuthorizationMethod) Selector() (string, error) {
	switch am {
	case AMApproveWithAuthorization:
		return string(constants.SelectorApproveWithAuthorization), nil
	case AMReceiveWithAuthorization:
		return string(constants.SelectorReceiveWithAuthorization), nil
	case AMReceiveWithAuthorizationEOA:
		return string(constants.SelectorReceiveWithAuthorizationEOA), nil
	}
	return string(constants.SelectorUndefined), errors.New("unknown authorization method")
}

type AuthorizationConfig struct {
	Domain              apitypes.TypedDataDomain
	AuthorizationMethod *AuthorizationMethod
}
