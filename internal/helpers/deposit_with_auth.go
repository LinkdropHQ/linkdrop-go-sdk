package helpers

import (
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk/internal/constants"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strings"
)

func GetDepositAuthorization(
	signTypedData types.SignTypedDataCallback,
	sender common.Address,
	to common.Address,
	amount *big.Int,
	validAfter int64,
	validBefore int64,
	transferId common.Address,
	expiration *big.Int,
	domain types.TypedDataDomain,
	token types.Token,
	feeAmount *big.Int,
	authSelector constants.Selector,
	authorizationMethod *types.AuthorizationMethod,
) ([]byte, error) {

	// If authorizationMethod is explicitly defined
	if authorizationMethod != nil {
		if *authorizationMethod == types.AMApproveWithAuthorization {
			return getDepositAuthorizationApprove(
				sender,
				to,
				amount,
				validAfter,
				validBefore,
				transferId,
				expiration,
				feeAmount,
				domain,
				signTypedData,
			)
		}
		return getDepositAuthorizationReceive(
			sender,
			to,
			amount,
			validAfter,
			validBefore,
			transferId,
			expiration,
			feeAmount,
			domain,
			authSelector,
			signTypedData,
		)
	}

	// Default behavior for Polygon chain
	if token.ChainId == types.ChainIdPolygon {
		if token.Address == constants.TAUsdcBridgedPolygon {
			return getDepositAuthorizationApprove(
				sender,
				to,
				amount,
				validAfter,
				validBefore,
				transferId,
				expiration,
				feeAmount,
				domain,
				signTypedData,
			)
		}
	}

	// Default fallback to ReceiveWithAuthorization
	return getDepositAuthorizationReceive(
		sender,
		to,
		amount,
		validAfter,
		validBefore,
		transferId,
		expiration,
		feeAmount,
		domain,
		authSelector,
		signTypedData,
	)
}

func getDepositAuthorizationApprove(
	sender common.Address,
	to common.Address,
	amount *big.Int,
	validAfter int64,
	validBefore int64,
	transferId common.Address,
	expiration *big.Int,
	feeAmount *big.Int,
	domain types.TypedDataDomain,
	signTypedData types.SignTypedDataCallback,
) ([]byte, error) {

	// Define the EIP-712 types
	t := map[string][]types.TypedDataField{
		"ApproveWithAuthorization": {
			{Type: "address", Name: "owner"},
			{Type: "address", Name: "spender"},
			{Type: "uint256", Name: "value"},
			{Type: "uint256", Name: "validAfter"},
			{Type: "uint256", Name: "validBefore"},
			{Type: "bytes32", Name: "nonce"},
		},
	}

	// Compute the nonce
	nonce := crypto.Keccak256Hash([]byte(fmt.Sprintf("%s%s%s%s", sender.Hex(), transferId.Hex(), amount.String(), expiration.String())))

	// Create the message
	message := map[string]interface{}{
		"owner":       sender,
		"spender":     to,
		"value":       amount,
		"validAfter":  validAfter,
		"validBefore": validBefore,
		"nonce":       nonce.Bytes(),
	}

	// Sign the typed data
	signature, err := signTypedData(domain, t, message)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to sign typed data: %w", err)
	}

	// Split the signature
	if len(signature) != 65 {
		return []byte{}, fmt.Errorf("invalid signature length: %d", len(signature))
	}
	r := signature[:32]
	s := signature[32:64]
	v := uint8(signature[64]) + 27 // Ensure Ethereum-specific "v" is set correctly

	// Encode the authorization using ABI encoding
	// Define the ABI definition for encoding
	abiJSON := `[{"name":"ApproveWithAuthorization","type":"function","inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"},{"name":"value","type":"uint256"},{"name":"validAfter","type":"uint256"},{"name":"validBefore","type":"uint256"},{"name":"nonce","type":"bytes32"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"}]}]`
	contractAbi, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return []byte{}, fmt.Errorf("failed to parse ABI JSON: %w", err)
	}

	// ABI encode the authorization
	authorizationData, err := contractAbi.Pack("ApproveWithAuthorization",
		sender, to, amount, big.NewInt(int64(validAfter)), big.NewInt(int64(validBefore)), nonce.Bytes(), v, r, s,
	)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to encode authorization data: %w", err)
	}

	return authorizationData, nil
}

func getDepositAuthorizationReceive(
	sender common.Address,
	to common.Address,
	amount *big.Int,
	validAfter int64,
	validBefore int64,
	transferId common.Address,
	expiration *big.Int,
	feeAmount *big.Int,
	domain types.TypedDataDomain,
	authSelector constants.Selector,
	signTypedData types.SignTypedDataCallback,
) ([]byte, error) {

	// Compute the nonce
	nonce := crypto.Keccak256Hash([]byte(fmt.Sprintf("%s%s%s%s%s", sender.Hex(), transferId.Hex(), amount.String(), expiration.String(), feeAmount.String())))

	// Create the message
	message := map[string]interface{}{
		"from":        sender,
		"to":          to,
		"value":       amount,
		"validAfter":  validAfter,
		"validBefore": validBefore,
		"nonce":       nonce.Bytes(),
	}

	// Define EIP-712 types for ReceiveWithAuthorization
	t := map[string][]types.TypedDataField{
		"ReceiveWithAuthorization": {
			{Name: "from", Type: "address"},
			{Name: "to", Type: "address"},
			{Name: "value", Type: "uint256"},
			{Name: "validAfter", Type: "uint256"},
			{Name: "validBefore", Type: "uint256"},
			{Name: "nonce", Type: "bytes32"},
		},
	}

	// Sign the typed data
	signature, err := signTypedData(domain, t, message)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to sign typed data: %w", err)
	}

	// Prepare for ABI encoding
	abiJSON := `[{"name":"ReceiveWithAuthorization","type":"function","inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"value","type":"uint256"},{"name":"validAfter","type":"uint256"},{"name":"validBefore","type":"uint256"},{"name":"nonce","type":"bytes32"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"},{"name":"signature","type":"bytes"}]}]`
	contractAbi, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return []byte{}, fmt.Errorf("failed to parse ABI JSON: %w", err)
	}

	if authSelector == constants.SelectorReceiveWithAuthorizationEOA {
		// Legacy format: split signature
		if len(signature) != 65 {
			return []byte{}, fmt.Errorf("invalid signature length: %d", len(signature))
		}
		r := signature[:32]
		s := signature[32:64]
		v := uint8(signature[64]) + 27 // Ethereum-specific adjustment for "v"

		// Encode using ABI
		authorizationData, err := contractAbi.Pack("ReceiveWithAuthorization",
			message["from"],
			message["to"],
			message["value"],
			message["validAfter"],
			message["validBefore"],
			message["nonce"],
			v,
			r,
			s,
		)
		if err != nil {
			return []byte{}, fmt.Errorf("failed to encode legacy authorization: %w", err)
		}
		return authorizationData, nil
	}
	// Modern format: include full signature as bytes
	authorizationData, err := contractAbi.Pack("ReceiveWithAuthorization",
		message["from"],
		message["to"],
		message["value"],
		message["validAfter"],
		message["validBefore"],
		message["nonce"],
		signature,
	)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to encode modern authorization: %w", err)
	}
	return authorizationData, nil
}
