package helpers

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
	"net/url"
	"strconv"
)

func defineSig(signatureLength int, signatureHex string) []byte {
	signature, err := base58.Decode(signatureHex)
	if err != nil {
		return []byte{}
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
	urlParts, err := url.ParseQuery(link)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query parameters: %w", err)
	}

	// Extract parameters with defaults if not present
	params := map[string]string{
		"linkKey":         urlParts.Get("k"),
		"signature":       urlParts.Get("sg"),
		"transferId":      urlParts.Get("i"),
		"chainId":         urlParts.Get("c"),
		"version":         urlParts.Get("v"),
		"signatureLength": urlParts.Get("sgl"),
		"encryptionKey":   urlParts.Get("m"),
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
	if params["transferId"] != "" {
		transferIdBytes, err := base58.Decode(params["transferId"])
		if err != nil {
			return nil, err
		}
		transferId = common.BytesToAddress(transferIdBytes)
	} else {
		// Derive wallet address from linkKey
		privateKey, err := crypto.ToECDSA(linkKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to derive private key from linkKey: %w", err)
		}
		transferId = crypto.PubkeyToAddress(privateKey.PublicKey)
	}

	// Parse chainId
	chainId, err := strconv.Atoi(params["chainId"])
	if err != nil {
		return nil, errors.New("invalid chainId value")
	}

	l := &types.Link{
		SenderSig:  senderSig,
		LinkKey:    linkKey,
		TransferId: transferId,
		ChainId:    types.ChainId(chainId),
		Version:    params["version"],
	}

	// Handle optional encryptionKey
	if params["encryptionKey"] != "" {
		var ek [32]byte
		ekDecoded, err := hex.DecodeString(params["encryptionKey"])
		if err != nil {
			return nil, errors.New("invalid encryptionKey value")
		}
		copy(ek[:], ekDecoded)
		l.EncryptionKeyLinkParam = &ek
	}

	return l, nil
}
