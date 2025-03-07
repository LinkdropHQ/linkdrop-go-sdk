package crypto

import (
	"encoding/hex"
	"errors"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/sha3"
)

const (
	NonceLength = 24
	KeyLength   = 32
	Type0       = 0
	TypeLength  = 1
)

func Keccak256(input string) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

// Converts a hex string to a byte array.
func fromHex(hexStr string) ([]byte, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Encodes the type as a single byte.
func encodeTypeByte(t byte) []byte {
	return []byte{t}
}

// Encrypt Encrypts a message using nacl.secretbox with TYPE_0 format:
// [type(1 byte), iv(24 bytes), sealed(...)] -> Hex encoded.
func Encrypt(
	message []byte,
	symKey [KeyLength]byte,
	iv [NonceLength]byte,
	getRandomBytes types.RandomBytesCallback,
) (encryptedMessage []byte, err error) {
	// Determine the IV (nonce)
	if iv == [NonceLength]byte{} {
		copy(iv[:], getRandomBytes(NonceLength))
	}

	// Encrypt the message with secretbox
	var naclKey [KeyLength]byte
	copy(naclKey[:], symKey[:])
	var naclNonce [NonceLength]byte
	copy(naclNonce[:], iv[:])

	sealed := secretbox.Seal(nil, message, &naclNonce, &naclKey)

	// Combine [type(1 byte), iv(24 bytes), sealed(...)]
	combined := append(encodeTypeByte(Type0), append(iv[:], sealed...)...)
	return combined, nil
}

// Decrypt decrypts a TYPE_0 formatted message using nacl.secretbox.
// Expects [type(1 byte), iv(24 bytes), sealed(...)] format.
//
// Inputs:
// - symKey: Hex-encoded 32-byte symmetric key
// - encoded: Hex-encoded encrypted message
//
// Output:
// - The decrypted message as a string, or an error in case of failure.
func Decrypt(encoded []byte, symKey [KeyLength]byte) (string, error) {
	// Decode the encrypted message from hex
	if len(encoded) < (TypeLength + NonceLength) {
		return "", errors.New("invalid encoded message format")
	}

	// Extract type byte (1st byte)
	typeByte := encoded[0]
	if typeByte != Type0 {
		return "", errors.New("invalid type byte, expected TYPE_0")
	}

	// Extract the IV (nonce) and sealed data
	ivStart := TypeLength
	ivEnd := ivStart + NonceLength
	iv := encoded[ivStart:ivEnd]
	if len(iv) != NonceLength {
		return "", errors.New("invalid IV length")
	}
	sealed := encoded[ivEnd:]

	// Prepare nacl key and nonce
	var naclKey [KeyLength]byte
	copy(naclKey[:], symKey[:])
	var naclNonce [NonceLength]byte
	copy(naclNonce[:], iv)

	// Decrypt the sealed message using secretbox
	decrypted, ok := secretbox.Open(nil, sealed, &naclNonce, &naclKey)
	if !ok {
		return "", errors.New("failed to decrypt")
	}

	// Return the plaintext message (as a string)
	return string(decrypted), nil
}
