package helpers

import (
	"errors"
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
	"net/url"
	"strconv"
	"strings"
)

func defineSig(signatureLength int, signatureHex string) []byte {
	signature, err := base58.Decode(signatureHex)
	if err != nil {
		return nil
	}

	paddedSignature := make([]byte, signatureLength)

	if len(signature) < signatureLength {
		// Pad with zeros if shorter
		offset := signatureLength - len(signature)
		copy(paddedSignature[offset:], signature)
	} else {
		signature = signature[:signatureLength]
	}
	return signature
}

func DecodeLink(link string) (*types.Link, error) {
	var queryArgs url.Values
	parsedLink, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	if strings.Contains(link, "/#/") {
		queryArgs, err = url.ParseQuery(strings.Split(parsedLink.Fragment, "?")[1])
	} else {
		queryArgs = parsedLink.Query()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse query parameters: %w", err)
	}

	// Extract parameters with defaults if not present
	params := map[string]string{
		"linkKey":         queryArgs.Get("k"),
		"signature":       queryArgs.Get("sg"),
		"transferId":      queryArgs.Get("i"),
		"chainId":         queryArgs.Get("c"),
		"version":         queryArgs.Get("v"),
		"signatureLength": queryArgs.Get("sgl"),
		"encryptionKey":   queryArgs.Get("m"),
	}

	// Set default values for version and signature length
	if params["version"] == "" {
		params["version"] = "1"
	}
	if params["signatureLength"] == "" {
		params["signatureLength"] = "65"
	}
	signatureLength, err := strconv.Atoi(params["signatureLength"])
	if err != nil {
		return nil, errors.New("invalid signature length value")
	}

	// Decode linkKey
	linkKeyBytes, err := base58.Decode(params["linkKey"])
	if err != nil {
		return nil, err
	}
	if linkKeyBytes == nil {
		return nil, errors.New("failed to decode linkKey")
	}
	linkKey, err := crypto.ToECDSA(linkKeyBytes)
	if err != nil {
		return nil, err
	}

	senderSig := defineSig(signatureLength, params["signature"])

	// Decode transferId or generate from linkKey
	var transferId common.Address
	if params["transferId"] == "" {
		// Derive wallet address from linkKey
		privateKey, err := crypto.ToECDSA(linkKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to derive private key from linkKey: %w", err)
		}
		transferId = crypto.PubkeyToAddress(privateKey.PublicKey)
	} else {
		transferIdBytes, err := base58.Decode(params["transferId"])
		if err != nil {
			return nil, err
		}
		transferId = common.BytesToAddress(transferIdBytes)
	}

	// Parse chainId
	chainId, err := strconv.Atoi(params["chainId"])
	if err != nil {
		return nil, errors.New("invalid chainId value")
	}

	l := &types.Link{
		SenderSignature: senderSig,
		LinkKey:         *linkKey,
		TransferId:      transferId,
		ChainId:         types.ChainId(chainId),
		Version:         params["version"],
	}

	// Handle optional encryptionKey
	if params["encryptionKey"] != "" {
		l.Message = &types.EncryptedMessage{
			Data:    nil,
			LinkKey: types.MessageLinkKey(params["encryptionKey"]),
		}
	}

	return l, nil
}
