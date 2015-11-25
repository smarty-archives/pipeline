package numeric

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
)

func Uint64ToString(value uint64) string {
	return strconv.FormatUint(value, Base10Encoding)
}

func StringToUint64(value string) uint64 {
	parsed, _ := strconv.ParseUint(value, Base10Encoding, 64)
	return parsed
}

func BinaryToUint64(bytes []byte) uint64 {
	var sum uint64

	length := len(bytes)

	if length > Uint64ByteLength {
		bytes = bytes[0:Uint64ByteLength]
		length = Uint64ByteLength
	}

	for i, b := range bytes {
		shift := uint64((length - i - 1) * 8)
		sum |= uint64(b) << shift
	}

	return sum
}

func GUIDToUint64(value string) uint64 {
	if len(value) != GUIDStringLength {
		return 0
	}

	value = value[0:8]
	if raw, err := hex.DecodeString(value); err != nil {
		return 0
	} else {
		return uint64(binary.BigEndian.Uint32(raw)) // .NET stores GUIDs in big endian
	}
}

func Uint32ToGUID(value uint32) string {
	var target [16]byte
	binary.BigEndian.PutUint32(target[:], value)
	return formatGUIDString(target)
}
func formatGUIDString(value [16]byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		value[0:4],
		value[4:6],
		value[6:8],
		value[8:10],
		value[10:16])
}

const (
	Base10Encoding   = 10
	Uint64ByteLength = 8
	GUIDStringLength = 36
)
