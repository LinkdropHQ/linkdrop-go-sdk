package linkdrop

import (
	"encoding/json"
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk/internal/helpers"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

type Client struct {
	config *Config
}

func (c *Client) RedeemRecoveredLink(apiHost, apiKey, receiver, sender, escrow, transferID, receiverSig, senderSig, token string) ([]byte, error) {
	body, _ := json.Marshal(map[string]string{
		"receiver":     receiver,
		"sender":       sender,
		"escrow":       escrow,
		"transfer_id":  transferID,
		"receiver_sig": receiverSig,
		"sender_sig":   senderSig,
		"token":        token,
	})
	return helpers.Request(fmt.Sprintf("%s/redeem-recovered", apiHost), "POST", helpers.DefineHeaders(apiKey), body)
}

func (c *Client) RedeemLink(apiHost, apiKey, receiver, transferID, receiverSig, token, sender, escrow string) ([]byte, error) {
	body, _ := json.Marshal(map[string]string{
		"receiver":     receiver,
		"sender":       sender,
		"escrow":       escrow,
		"transfer_id":  transferID,
		"receiver_sig": receiverSig,
		"token":        token,
	})
	return helpers.Request(fmt.Sprintf("%s/redeem", apiHost), "POST", helpers.DefineHeaders(apiKey), body)
}

func (c *Client) GetTransferStatus(apiHost, apiKey, transferID string) ([]byte, error) {
	return helpers.Request(fmt.Sprintf("%s/payment-status/transfer/%s", apiHost, transferID), "GET", helpers.DefineHeaders(apiKey), nil)
}

func (c *Client) GetTransferStatusByTxHash(apiHost, apiKey, txHash string) ([]byte, error) {
	return helpers.Request(fmt.Sprintf("%s/payment-status/transaction/%s", apiHost, txHash), "GET", helpers.DefineHeaders(apiKey), nil)
}

func (c *Client) GetFee(apiHost, apiKey, tokenAddress, sender, tokenType, transferID, expiration, amount, tokenID string) ([]byte, error) {
	query := helpers.CreateQueryString(map[string]string{
		"amount":        amount,
		"token_address": tokenAddress,
		"sender":        sender,
		"token_type":    tokenType,
		"transfer_id":   transferID,
		"expiration":    expiration,
		"token_id":      tokenID,
	})
	return helpers.Request(fmt.Sprintf("%s/fee?%s", apiHost, query), "GET", helpers.DefineHeaders(apiKey), nil)
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
	apiHost, apiKey, token, tokenType, sender, escrow, transferID, expiration, txHash, feeAuthorization string,
	amount, feeAmount, totalAmount string,
	feeToken, encryptedSenderMessage string,
) ([]byte, error) {
	body, err := json.Marshal(map[string]string{
		"sender":                   sender,
		"escrow":                   escrow,
		"transfer_id":              transferID,
		"token":                    token,
		"token_type":               tokenType,
		"expiration":               expiration,
		"tx_hash":                  txHash,
		"fee_authorization":        feeAuthorization,
		"amount":                   amount,
		"fee_amount":               feeAmount,
		"total_amount":             totalAmount,
		"fee_token":                feeToken,
		"encrypted_sender_message": encryptedSenderMessage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	return helpers.Request(fmt.Sprintf("%s/deposit", apiHost), "POST", helpers.DefineHeaders(apiKey), body)
}

func (c *Client) DepositERC721(
	apiHost, apiKey, token, tokenType, sender, escrow, transferID, expiration, txHash, feeAuthorization, tokenID string,
	feeAmount, totalAmount, feeToken, encryptedSenderMessage string,
) ([]byte, error) {
	// Create the body of the request
	body, err := json.Marshal(map[string]string{
		"sender":                   sender,
		"escrow":                   escrow,
		"transfer_id":              transferID,
		"token":                    token,
		"token_type":               tokenType,
		"expiration":               expiration,
		"tx_hash":                  txHash,
		"fee_authorization":        feeAuthorization,
		"token_id":                 tokenID,
		"amount":                   "1", // ERC721 always has a fixed quantity of 1
		"fee_amount":               feeAmount,
		"total_amount":             totalAmount,
		"fee_token":                feeToken,
		"encrypted_sender_message": encryptedSenderMessage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	return helpers.Request(fmt.Sprintf("%s/deposit-erc721", apiHost), "POST", helpers.DefineHeaders(apiKey), body)
}

func (c *Client) DepositERC1155(
	apiHost, apiKey, token, tokenType, sender, escrow, transferID, expiration, txHash, feeAuthorization, tokenID, amount string,
	feeAmount, totalAmount, feeToken, encryptedSenderMessage string,
) ([]byte, error) {
	// Create the body of the request
	body, err := json.Marshal(map[string]string{
		"sender":                   sender,
		"escrow":                   escrow,
		"transfer_id":              transferID,
		"token":                    token,
		"token_type":               tokenType,
		"expiration":               expiration,
		"tx_hash":                  txHash,
		"fee_authorization":        feeAuthorization,
		"token_id":                 tokenID,
		"amount":                   amount,
		"fee_amount":               feeAmount,
		"total_amount":             totalAmount,
		"fee_token":                feeToken,
		"encrypted_sender_message": encryptedSenderMessage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	return helpers.Request(fmt.Sprintf("%s/deposit-erc1155", apiHost), "POST", helpers.DefineHeaders(apiKey), body)
}

func (c *Client) DepositWithAuthorization(
	apiHost, apiKey, token, tokenType, sender, escrow, transferID, expiration, authorization, authorizationSelector, feeAuthorization, amount string,
	feeAmount, totalAmount, encryptedSenderMessage string,
) ([]byte, error) {
	// Create the body of the request
	body, err := json.Marshal(map[string]string{
		"sender":                   sender,
		"escrow":                   escrow,
		"transfer_id":              transferID,
		"token":                    token,
		"token_type":               tokenType,
		"expiration":               expiration,
		"amount":                   amount,
		"authorization":            authorization,
		"authorization_selector":   authorizationSelector,
		"fee_amount":               feeAmount,
		"total_amount":             totalAmount,
		"fee_authorization":        feeAuthorization,
		"encrypted_sender_message": encryptedSenderMessage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	return helpers.Request(fmt.Sprintf("%s/deposit-with-authorization", apiHost), "POST", helpers.DefineHeaders(apiKey), body)
}
