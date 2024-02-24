package bitbuffer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	t.Run("", func(t *testing.T) {
		bb := NewBitBuffer([]uint8{0x45, 0x54})

		x, _ := bb.ReadUint8()
		y, _ := bb.ReadUint8()
		assert.Equal(t, uint8(0x45), x)
		assert.Equal(t, uint8(0x54), y)
	})

	t.Run("", func(t *testing.T) {
		bb := NewBitBuffer([]uint8{0b1111_1111})

		x, _ := bb.ReadBits(3)
		y, _ := bb.ReadBits(1)
		z, _ := bb.ReadBits(4)
		assert.Equal(t, uint64(0b111), x)
		assert.Equal(t, uint64(0b1), y)
		assert.Equal(t, uint64(0b1111), z)
	})

	t.Run("", func(t *testing.T) {
		bb := NewBitBuffer([]uint8{0b0100_0101})

		x, _ := bb.ReadBits(3)
		y, _ := bb.ReadBits(1)
		z, _ := bb.ReadBits(4)
		assert.Equal(t, uint64(0b101), x)
		assert.Equal(t, uint64(0b0), y)
		assert.Equal(t, uint64(0b0100), z)
	})
}
