package helpers

import (
	"encoding/hex"
	"fmt"
)

// NumberToHexString converts an integer to a hexadecimal string
func NumberToHexString(number int64) string {
	return fmt.Sprintf("%02X", number)
}

func ToHex(data []byte) string {
	return hex.EncodeToString(data)
}
