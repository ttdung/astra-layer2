package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// FwdcommitmentKeyPrefix is the prefix to retrieve all Fwdcommitment
	FwdcommitmentKeyPrefix = "Fwdcommitment/value/"
)

// FwdcommitmentKey returns the store key to retrieve a Fwdcommitment from the index fields
func FwdcommitmentKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
