package helpers

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

const EIP191Prefix = "\x19Ethereum Signed Message:\n"

// GenerateReceiverSig takes a private key in hexadecimal format (linkKey) and a receiver Ethereum address,
// and generates a signature for the receiver address using the Keccak256 hashing algorithm and the private key.
//
// Parameters:
// - linkKey: A string representing the private key in hexadecimal format.
// - receiver: An Ethereum address represented by the `common.Address` type.
//
// Returns:
// - A byte slice containing the generated signature.
// - An error if the private key is invalid or signing fails.
//
// Example:
//
//	  linkKey := "your-private-key-hex"
//	  receiver := common.HexToAddress("0xReceiverAddress")
//	  sig, err := GenerateReceiverSig(linkKey, receiver)
//	  if err != nil {
//		   // Handle error
//	  }
//	  fmt.Printf("Generated signature: %x\n", sig)
func GenerateReceiverSig(linkKey *ecdsa.PrivateKey, receiver common.Address) ([]byte, error) {
	messageHash := sha3.NewLegacyKeccak256()
	messageHash.Write(receiver.Bytes())
	message := messageHash.Sum(nil)

	prefix := []byte(EIP191Prefix)
	length := []byte(fmt.Sprintf("%d", len(message)))

	eip191Hash := sha3.NewLegacyKeccak256()
	eip191Hash.Write(append(append(prefix, length...), message...))

	signature, err := crypto.Sign(eip191Hash.Sum(nil), linkKey)
	if signature[64] < 27 {
		signature[64] += 27
	}
	return signature, err
}
