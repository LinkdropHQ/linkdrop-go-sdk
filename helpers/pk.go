package helpers

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func AddressFromPrivateKey(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, fmt.Errorf("failed to type assert public key to *ecdsa.PublicKey")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA), nil
}

func PrivateKey(
	getRandomBytes types.RandomBytesCallback,
) (*ecdsa.PrivateKey, error) {
	seed := getRandomBytes(32)
	if len(seed) != 32 {
		return nil, fmt.Errorf("seed must be exactly 32 bytes")
	}
	return PrivateKeyFromHash(common.BytesToHash(seed))
}

func PrivateKeyFromHash(privateKeyBytes common.Hash) (*ecdsa.PrivateKey, error) {
	d := new(big.Int).SetBytes(privateKeyBytes.Bytes())
	curve := crypto.S256()
	order := curve.Params().N
	if d.Sign() <= 0 || d.Cmp(order) >= 0 {
		return nil, fmt.Errorf("seed is out of valid range for secp256k1 private key")
	}
	privateKey := new(ecdsa.PrivateKey)
	privateKey.D = d
	privateKey.PublicKey.Curve = curve
	privateKey.PublicKey.X, privateKey.PublicKey.Y = curve.ScalarBaseMult(privateKeyBytes.Bytes())

	return privateKey, nil
}
