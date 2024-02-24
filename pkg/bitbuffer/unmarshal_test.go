package bitbuffer

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type StructWithUnmarshaler struct {
	A bool
	B bool
}

func (foo *StructWithUnmarshaler) Unmarshal(bb *BitBuffer, _ reflect.Value, _ reflect.StructTag) error {
	a, err := bb.ReadUint8()
	if err != nil {
		return err
	}
	b, err := bb.ReadUint8()
	if err != nil {
		return err
	}
	foo.A = a != 0
	foo.B = b != 0
	return nil
}

func TestBitBuffer_Unmarshal(t *testing.T) {
	t.Run("unmarshalling 2 bools succeeds", func(t *testing.T) {
		type Foo struct {
			A bool
			B bool
		}

		data := []uint8{0b10}

		var foo Foo
		err := Unmarshal(data, &foo)
		assert.Equal(t, nil, err)
		assert.Equal(t, Foo{A: false, B: true}, foo)
	})

	t.Run("unmarshalling a field with custom unmarshaler succeeds", func(t *testing.T) {
		data := []uint8{0, 1}

		var s StructWithUnmarshaler
		err := Unmarshal(data, &s)
		assert.Equal(t, nil, err)
		assert.Equal(t, StructWithUnmarshaler{A: false, B: true}, s)
	})

	t.Run("fully unmarshalling a struct withe exact data succeeds", func(t *testing.T) {
		type Foo struct {
			A uint8
			B uint8
		}

		data := []uint8{0x55, 0x90}

		var foo Foo
		err := UnmarshalFully(data, &foo)
		assert.Equal(t, nil, err)
		assert.Equal(t, Foo{A: 0x55, B: 0x90}, foo)
	})

	t.Run("fully unmarshalling a struct with extra data fails", func(t *testing.T) {
		type Foo struct {
			A uint8
			B uint8
		}

		data := []uint8{0x55, 0x90, 0x88}

		var foo Foo
		err := UnmarshalFully(data, &foo)
		assert.EqualError(t, err, "8 bits remaining in bitbuffer after unmarshal")
	})

	t.Run("values that are not byte-aligned cannot be unmarshalled fully", func(t *testing.T) {
		type Foo struct {
			A bool
			B bool
		}

		data := []uint8{0b10}

		var foo Foo
		err := UnmarshalFully(data, &foo)
		assert.EqualError(t, err, "6 bits remaining in bitbuffer after unmarshal")
	})
}
