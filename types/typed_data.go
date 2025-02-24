package types

import "github.com/ethereum/go-ethereum/common"

// TypedDataDomain represents the domain for signing data
type TypedDataDomain struct {
	Name              string
	Version           string
	ChainId           ChainId
	VerifyingContract *common.Address
	Salt              *[]byte
}

// TypedDataField represents a field inside the data types for signing
type TypedDataField struct {
	Name string
	Type string
}
