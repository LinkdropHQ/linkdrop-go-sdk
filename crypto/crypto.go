package crypto

import (
	"errors"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	NonceLength = 24
	KeyLength   = 32
	Type0       = 0
	TypeLength  = 1
)

// Encodes the type as a single byte.
func encodeTypeByte(t byte) []byte {
	return []byte{t}
}

// Encrypt Encrypts a message using nacl.secretbox with TYPE_0 format:
// [type(1 byte), iv(24 bytes), sealed(...)] -> Hex encoded.
// NOTE: iv can be provided as a random [NonceLength]byte array
func Encrypt(
	message []byte,
	symKey [KeyLength]byte,
	iv [NonceLength]byte,
) (encryptedMessage []byte, err error) {
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
