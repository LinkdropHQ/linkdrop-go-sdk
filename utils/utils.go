package utils

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/LinkdropHQ/linkdrop-go-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	geth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"log"
	"math/big"
)

func GetRandomBytes(length int64) []byte {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate random bytes: %v", err)
	}
	return b
}

func SignTypedData(typedData apitypes.TypedData, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	typedDataHash, _, err := apitypes.TypedDataAndHash(typedData)
	if err != nil {
		return nil, fmt.Errorf("failed to hash typed data: %w", err)
	}
	signature, err := crypto.Sign(typedDataHash[:], privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign the typed data: %w", err)
	}
	signature[64] += 27
	return signature, nil
}

func SendTransaction(
	chainId *big.Int,
	to common.Address,
	value *big.Int,
	data []byte,

	client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
) (transaction *types.Transaction, err error) {
	rpcChainId, err := client.NetworkID(context.Background())
	if err != nil {
		return
	}
	if rpcChainId.Cmp(chainId) != 0 {
		return nil, fmt.Errorf("wrong RPC chain ID: %s, expected: %s", rpcChainId, chainId)
	}

	sender, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return
	}

	nonce, err := client.PendingNonceAt(context.Background(), sender.From)
	if err != nil {
		return
	}
	log.Println("Nonce:", nonce)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  sender.From,
		To:    &to,
		Value: value,
		Data:  data,
	})
	if err != nil {
		gasLimit = 100000
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	gasPrice.Add(gasPrice, gasPrice)
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}

	tx := geth_types.NewTx(&geth_types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	signedTx, err := geth_types.SignTx(tx, geth_types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		return
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return
	}

	transaction = &types.Transaction{
		Hash: signedTx.Hash(),
		Type: "tx",
	}
	return
}
