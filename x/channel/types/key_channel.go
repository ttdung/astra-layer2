package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ChannelKeyPrefix is the prefix to retrieve all Channel
	ChannelKeyPrefix = "Channel/value/"
)

// ChannelKey returns the store key to retrieve a Channel from the index fields
func ChannelKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
