package helpers

import "fmt"

// NumberToHexString converts an integer to a hexadecimal string
func NumberToHexString(number int64) string {
	return fmt.Sprintf("%02X", number)
}
