package helpers

import (
	"crypto/sha256"
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk/internal/crypto"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mr-tron/base58"
)

func EncryptMessage(
	message string,
	transferID common.Address,
	chainID types.ChainId,
	encryptionKeyLength int64,
	getRandomBytes types.RandomBytesCallback,
	signTypedData types.SignTypedDataCallback,
) (*types.EncryptedMessage, error) {
	encryptionKey, encryptionKeyLinkParam, err := CreateMessageEncryptionKey(transferID.Hex(), signTypedData, chainID, encryptionKeyLength)
	if err != nil {
		return nil, err
	}

	encryptedSenderMessage, err := crypto.Encrypt([]byte(message), encryptionKey, []byte{}, getRandomBytes)
	if err != nil {
		return nil, err
	}

	// Convert encryption key length to hexadecimal
	encryptionKeyLengthAsHex := NumberToHexString(encryptionKeyLength)

	// Build the result
	result := &types.EncryptedMessage{
		Message:       fmt.Sprintf("%s%s", encryptionKeyLengthAsHex, encryptedSenderMessage),
		EncryptionKey: encryptionKeyLinkParam,
	}
	return result, nil
}

// CreateMessageEncryptionKey creates a message encryption key
func CreateMessageEncryptionKey(
	transferID string,
	signTypedData types.SignTypedDataCallback,
	chainId types.ChainId,
	encryptionKeyLength int64,
) (encryptionKey []byte, encryptionKeyLinkParam []byte, err error) {
	domain := types.TypedDataDomain{
		Name:    "MyEncryptionScheme",
		Version: "1",
		ChainId: chainId,
	}

	typedData := map[string][]types.TypedDataField{
		"EncryptionMessage": {
			{Name: "seed", Type: "string"},
		},
	}

	seed := fmt.Sprintf("Encrypting message (transferId: %s)", transferID)
	value := map[string]interface{}{
		"seed": seed,
	}

	// Generating signature
	signature, err := signTypedData(domain, typedData, value)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	// Calculating the initial encryption key
	encryptionKeyInitial := sha256.Sum256([]byte(signature))

	// Encoding to Base58 and trimming
	encryptionKeyModified := base58.Encode(encryptionKeyInitial[:])
	if int64(len(encryptionKeyModified)) > encryptionKeyLength {
		encryptionKeyModified = encryptionKeyModified[:encryptionKeyLength]
	}

	// Converting Base58 code to byte array before hashing
	encryptionKeyModifiedBytes, err := base58.Decode(encryptionKeyModified)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	// Final encryption
	finalHash := sha256.Sum256(encryptionKeyModifiedBytes)

	return finalHash[:], encryptionKeyInitial[:], nil
}
