package numeric

import (
	"crypto/rand"
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

func NewGUID() string {
	buffer := make([]byte, GUIDByteLength)
	rand.Read(buffer)

	encoded := hex.EncodeToString(buffer)
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		encoded[:8],
		encoded[8:12],
		encoded[12:16],
		encoded[16:20],
		encoded[20:])
}

func GUIDToString(value string) uint64 {
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

const Base10Encoding = 10
const Uint64ByteLength = 8
const GUIDByteLength = 16
const GUIDStringLength = 36
