package numeric

import (
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func Test(t *testing.T) {
	assert := assertions.New(t)

	assert.So(Uint64ToString(uint64(maxUint64)), should.Equal, maxUint64String)

	assert.So(StringToUint64(maxUint64String), should.Equal, maxUint64)
	assert.So(StringToUint64("blah"), should.Equal, 0)

	assert.So(BinaryToUint64([]byte{0, 0, 0, 42}), should.Equal, 42)
	assert.So(BinaryToUint64([]byte{255, 255, 255, 255}), should.Equal, maxInt)
	assert.So(BinaryToUint64([]byte{255, 255, 255, 255, 255, 255, 255, 255}), should.Equal, maxUint64)
	assert.So(BinaryToUint64([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255}), should.Equal, maxUint64) // only convert up to 8 bytes

	assert.So(GUIDToUint64(guidString), should.Equal, guidUint64)
	assert.So(GUIDToUint64("wrong length"), should.Equal, 0)
	assert.So(GUIDToUint64("malforms-0000-0000-0000-000000000000"), should.Equal, 0)
}

const (
	maxUint64       uint64 = 18446744073709551615
	maxUint64String string = "18446744073709551615"
	maxInt          int    = 4294967295
	guidString      string = "2c8991cf-0000-0000-0000-000000000000"
	guidUint64      uint64 = 747213263
)
