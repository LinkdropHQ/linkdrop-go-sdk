package linkdrop

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk/internal/helpers"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Client struct {
	config *Config
}

// RedeemRecoveredLink allows a receiver to redeem a link that has been recovered.
// This function works similarly to RedeemLink but is used in cases where both the
// receiver and sender signatures are required for recovery.
//
// Parameters:
// - receiver: The Ethereum address of the receiver.
// - transferId: The unique identifier of the transfer.
// - receiverSig: The signature of the receiver as a byte slice.
// - senderSig: The signature of the sender as a hex string.
// - sender: The Ethereum address of the sender.
// - escrow: The Ethereum address of the escrow account.
// - token: The token details (address and chain ID).
//
// Returns:
// - []byte: A JSON-encoded response from the API upon successful redemption.
// - error: An error object if the redemption fails, otherwise nil.
//
// Notes:
// - This function sends a POST request to the API endpoint `/redeem-recovered`.
// - Ensure all required parameters are valid and signatures are properly formed.
// - The API validates sender, receiver, and escrow information against the signatures.
func (c *Client) RedeemRecoveredLink(
	receiver common.Address,
	transferId common.Address,
	receiverSig []byte,
	senderSig []byte,
	sender common.Address,
	escrow common.Address,
	token types.Token,
) ([]byte, error) {
	body, _ := json.Marshal(map[string]string{
		"receiver":     receiver.Hex(),
		"sender":       sender.Hex(),
		"escrow":       escrow.Hex(),
		"transfer_id":  transferId.Hex(),
		"receiver_sig": "0x" + hex.EncodeToString(receiverSig),
		"sender_sig":   "0x" + hex.EncodeToString(senderSig),
		"token":        token.Address.Hex(),
	})
	return helpers.Request(fmt.Sprintf("%s/redeem-recovered", c.config.apiURL), "POST", helpers.DefineHeaders(c.config.apiKey), body)
}

// RedeemLink allows a receiver to redeem a link by providing details such as transfer ID,
// receiver signature, and optionally the sender, escrow, and token information.
//
// Parameters:
// - receiver: The Ethereum address of the receiver.
// - transferId: The unique identifier of the transfer.
// - receiverSig: The signature of the receiver as a byte slice.
// - sender: (Optional) The Ethereum address of the sender.
// - escrow: (Optional) The Ethereum address of the escrow account.
// - token: (Optional) The token details (address and chain ID).
//
// Returns:
// - []byte: A JSON-encoded response from the API upon successful redemption.
// - error: An error object if the redemption fails, otherwise nil.
//
// Notes:
// - This function sends a POST request to the API to redeem the link.
// - Ensure all required parameters are valid before calling this function.
// - If optional parameters (sender, escrow, or token) are not provided, they will be ignored in the request body.
func (c *Client) RedeemLink(
	chainId types.ChainId,
	receiver common.Address,
	transferId common.Address,
	receiverSig []byte,
	sender *common.Address,
	escrow *common.Address,
	token *types.Token,
) ([]byte, error) {
	apiHost, err := helpers.DefineApiHost(c.config.apiURL, int64(chainId))
	if err != nil {
		return []byte{}, err
	}
	bodyRaw := map[string]string{
		"receiver":     receiver.Hex(),
		"transfer_id":  transferId.Hex(),
		"receiver_sig": "0x" + hex.EncodeToString(receiverSig),
	}
	if sender != nil {
		bodyRaw["sender"] = sender.Hex()
	}
	if escrow != nil {
		bodyRaw["escrow"] = escrow.Hex()
	}
	if token != nil {
		bodyRaw["token"] = token.Address.Hex()
	}
	body, _ := json.Marshal(bodyRaw)
	return helpers.Request(fmt.Sprintf("%s/redeem", apiHost), "POST", helpers.DefineHeaders(c.config.apiKey), body)
}

// GetTransferStatus retrieves the payment status of a transfer using its unique transfer ID.
//
// Parameters:
// - transferId: The unique identifier of the transfer.
//
// Returns:
// - []byte: A JSON-encoded response from the API containing the payment status.
// - error: An error object if the operation fails, otherwise nil.
//
// Notes:
// - The function sends a GET request to the API to fetch the transfer's payment status.
// - Ensure that the transfer ID is valid and corresponds to an existing transfer.
func (c *Client) GetTransferStatus(
	chainId types.ChainId,
	transferId common.Address,
) ([]byte, error) {
	apiHost, err := helpers.DefineApiHost(c.config.apiURL, int64(chainId))
	if err != nil {
		return []byte{}, err
	}
	return helpers.Request(fmt.Sprintf("%s/payment-status/transfer/%s", apiHost, transferId.Hex()), "GET", helpers.DefineHeaders(c.config.apiKey), nil)
}

// GetTransferStatusByTxHash retrieves the payment status of a transfer using its transaction hash.
//
// Parameters:
// - txHash: The transaction hash associated with the transfer.
//
// Returns:
// - []byte: A JSON-encoded response from the API containing the payment status.
// - error: An error object if the operation fails, otherwise nil.
//
// Notes:
// - This function sends a GET request to the API to fetch the transfer's payment status by transaction hash.
// - Ensure the transaction hash corresponds to a valid transfer and has been processed.
func (c *Client) GetTransferStatusByTxHash(txHash string) ([]byte, error) {
	return helpers.Request(fmt.Sprintf("%s/payment-status/transaction/%s", c.config.apiURL, txHash), "GET", helpers.DefineHeaders(c.config.apiKey), nil)
}

// GetFee calculates the transaction fee required for a transfer based on token details, sender's address, transfer ID,
// expiration time, and transfer amount.
//
// Parameters:
// - token: The token object consisting of the token's address, type, and ID.
// - sender: The Ethereum address of the sender.
// - transferId: The unique identifier of the transfer.
// - expiration: A pointer to a big.Int object representing the expiration time of the transfer (in seconds, typically).
// - amount: A pointer to a big.Int object representing the transfer amount.
//
// Returns:
// - []byte: A JSON-encoded response from the API containing the calculated fee details.
// - error: An error object if the operation fails, otherwise nil.
//
// Notes:
// - The function validates that `amount` and `expiration` are not nil, as they are required parameters.
// - If the token includes an ID, it is included in the request query string; otherwise, "0" is set as the default token ID.
// - The `CreateQueryString` helper function is used to construct the query string with the required parameters.
// - The function sends a GET request to the API, formatted as `apiURL/fee`, with the constructed query string.
// - Ensure the provided parameters are valid and match the expected values in the API.
func (c *Client) GetFee(
	token types.Token,
	sender common.Address,
	transferId common.Address,
	expiration *big.Int,
	amount *big.Int,
) ([]byte, error) {
	if amount == nil || expiration == nil {
		return nil, fmt.Errorf("amount and expiration are required")
	}
	tokenId := "0"
	if token.Id != nil {
		tokenId = token.Id.String()
	}
	query := helpers.CreateQueryString(map[string]string{
		"amount":        amount.String(),
		"token_address": token.Address.Hex(),
		"sender":        sender.Hex(),
		"token_type":    string(token.Type),
		"transfer_id":   transferId.Hex(),
		"expiration":    expiration.String(),
		"token_id":      tokenId,
	})
	return helpers.Request(fmt.Sprintf("%s/fee?%s", c.config.apiURL, query), "GET", helpers.DefineHeaders(c.config.apiKey), nil)
}

// GetHistory fetches the history of transfers related to a token and sender's address.
// It allows filtering for only active transfers and supports pagination with offset and limit.
//
// Parameters:
// - token: The token object containing the token's address and chain ID.
// - sender: The sender's Ethereum address.
// - onlyActive: A boolean flag to filter only active transfers (true) or all transfers (false).
// - offset: The number of records to skip for pagination.
// - limit: The maximum number of records to fetch for pagination.
//
// Returns:
// - []byte: A JSON-encoded response from the API containing the transfer history.
// - error: An error object if the operation fails, otherwise nil.
//
// Notes:
// - The function dynamically determines the API host based on the token's chain ID.
func (c *Client) GetHistory(token types.Token, sender common.Address, onlyActive bool, offset, limit int64) ([]byte, error) {
	apiHost, err := helpers.DefineApiHost(c.config.apiURL, int64(token.ChainId))
	if err != nil {
		return nil, err
	}
	query := helpers.CreateQueryString(map[string]string{
		"only_active":   fmt.Sprintf("%t", onlyActive),
		"offset":        fmt.Sprintf("%d", offset),
		"limit":         fmt.Sprintf("%d", limit),
		"token_address": token.Address.Hex(),
	})
	return helpers.Request(
		fmt.Sprintf("%s/payment-status/sender/%s/get-sender-history?%s", apiHost, sender.Hex(), query),
		"GET",
		helpers.DefineHeaders(c.config.apiKey),
		nil,
	)
}

// GetLimits fetches the limits of a specific token from the API.
//
// Parameters:
// - token: The token object containing the token's address, type, and chain ID.
//
// Returns:
// - []byte: A JSON-encoded response from the API containing the token limits.
// - error: An error object if the operation fails, otherwise nil.
//
// Notes:
// - The function determines the appropriate API host based on the token's chain ID.
func (c *Client) GetLimits(token types.Token) ([]byte, error) {
	apiHost, err := helpers.DefineApiHost(c.config.apiURL, int64(token.ChainId))
	if err != nil {
		return nil, err
	}
	query := helpers.CreateQueryString(map[string]string{
		"token_address": token.Address.Hex(),
		"token_type":    string(token.Type),
	})
	return helpers.Request(
		fmt.Sprintf("%s/limits?%s", apiHost, query),
		"GET",
		helpers.DefineHeaders(c.config.apiKey),
		nil,
	)
}

func (c *Client) Deposit(
	token types.Token,
	sender common.Address,
	escrow common.Address,
	transferId common.Address,
	expiration *big.Int,
	transaction *types.Transaction,
	fee types.CLFee,
	amount *big.Int,
	totalAmount *big.Int,
	encryptedSenderMessage string,
) ([]byte, error) {
	bodyRaw := map[string]string{
		"sender":                   sender.Hex(),
		"escrow":                   escrow.Hex(),
		"transfer_id":              transferId.Hex(),
		"token":                    token.Address.Hex(),
		"token_type":               string(token.Type),
		"expiration":               expiration.String(),
		"tx_hash":                  transaction.Hash.Hex(),
		"type":                     string(transaction.Type),
		"fee_authorization":        "0x" + fee.Authorization,
		"amount":                   amount.String(),
		"fee_amount":               fee.Amount.String(),
		"total_amount":             totalAmount.String(),
		"fee_token":                fee.Token.Address.Hex(),
		"encrypted_sender_message": encryptedSenderMessage,
		"token_id":                 token.Id.String(),
	}

	var endpoint string
	switch token.Type {
	case types.TokenTypeERC721:
		bodyRaw["amount"] = "1"
		bodyRaw["token_id"] = token.Id.String()
		endpoint = "%s/deposit-erc721"
	case types.TokenTypeERC1155:
		bodyRaw["token_id"] = token.Id.String()
		endpoint = "%s/deposit-erc1155"
	default:
		endpoint = "%s/deposit"
	}

	body, err := json.Marshal(bodyRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	return helpers.Request(fmt.Sprintf(endpoint, c.config.apiURL), "POST", helpers.DefineHeaders(c.config.apiKey), body)
}

func (c *Client) DepositWithAuthorization(
	token types.Token,
	sender common.Address,
	escrow common.Address,
	transferId common.Address,
	expiration *big.Int,
	authorization []byte,
	authorizationSelector string,
	fee types.CLFee,
	amount *big.Int,
	totalAmount *big.Int,
	encryptedSenderMessage string,
) ([]byte, error) {
	// Create the body of the request
	body, err := json.Marshal(map[string]string{
		"sender":                   sender.Hex(),
		"escrow":                   escrow.Hex(),
		"transfer_id":              transferId.Hex(),
		"token":                    token.Address.Hex(),
		"token_type":               string(token.Type),
		"expiration":               expiration.String(),
		"amount":                   amount.String(),
		"authorization":            "0x" + hex.EncodeToString(authorization),
		"authorization_selector":   authorizationSelector,
		"fee_amount":               fee.Amount.String(),
		"total_amount":             totalAmount.String(),
		"fee_authorization":        "0x" + fee.Authorization,
		"encrypted_sender_message": encryptedSenderMessage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	return helpers.Request(fmt.Sprintf("%s/deposit-with-authorization", c.config.apiURL), "POST", helpers.DefineHeaders(c.config.apiKey), body)
}
