package bitbuffer

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBitBuffer_WidthTag(t *testing.T) {
	t.Run("unmarshalling 2 bools with overridden widths succeeds", func(t *testing.T) {
		type Foo struct {
			A bool `bbwidth:"8"`
			B bool `bbwidth:"8"`
		}

		data := []uint8{0, 1}

		var foo Foo
		err := UnmarshalFully(data, &foo)
		assert.Equal(t, nil, err)
		assert.Equal(t, Foo{A: false, B: true}, foo)
	})

	t.Run("unmarshalling a field with a bad width tag fails", func(t *testing.T) {
		type StructWithBadWidthTag struct {
			A uint8
			B uint8 `bbwidth:"width"`
		}

		data := []uint8{12, 14}

		var s StructWithBadWidthTag
		err := Unmarshal(data, &s)
		assert.EqualError(t, err, "strconv.ParseUint: parsing \"width\": invalid syntax")
	})

	t.Run("unmarshalling a field with a zero width tag fail", func(t *testing.T) {
		type StructWithZeroWidthTag struct {
			A uint8
			B uint8 `bbwidth:"0"`
			C uint8
		}

		data := []uint8{12, 14}

		var s StructWithZeroWidthTag
		err := Unmarshal(data, &s)
		assert.EqualError(t, err, "bit width must be greater than zero")
	})

	t.Run("unmarshalling a field with an oversized width tag succeeds", func(t *testing.T) {
		type StructWithOversizedWidthTag struct {
			A uint8 `bbwidth:"65"`
			B uint8 `bbwidth:"7"`
			C uint8
		}

		data := []uint8{11, 22, 33, 44, 55, 66, 77, 88, 99, 111}

		var s StructWithOversizedWidthTag
		err := UnmarshalFully(data, &s)
		assert.Equal(t, nil, err)
		assert.Equal(t, StructWithOversizedWidthTag{A: 11, B: 99 >> 1, C: 111}, s)
	})
}

type ValidatedField uint8

func (vf ValidatedField) Validate() error {
	if vf > 8 {
		return errors.New("invalid")
	}
	return nil
}

type StructWithValidatedField struct {
	A uint8
	B ValidatedField `bbvalidate:"true"`
}

func TestBitBuffer_ValidateTag(t *testing.T) {
	t.Run("parsing a validated field with valid data succeeds", func(t *testing.T) {
		data := []uint8{12, 6}

		var s StructWithValidatedField
		err := Unmarshal(data, &s)
		assert.Equal(t, nil, err)
		assert.Equal(t, StructWithValidatedField{A: 12, B: 6}, s)
	})

	t.Run("parsing a validated field with invalid data fails", func(t *testing.T) {
		data := []uint8{12, 14}

		var s StructWithValidatedField
		err := Unmarshal(data, &s)
		assert.EqualError(t, err, "invalid")
	})

	t.Run("parsing a validated field with a bad tag fails", func(t *testing.T) {
		type StructWithBadValidateTag struct {
			A uint8
			B ValidatedField `bbvalidate:"validate"`
		}

		data := []uint8{12, 14}

		var s StructWithBadValidateTag
		err := Unmarshal(data, &s)
		assert.EqualError(t, err, "strconv.ParseBool: parsing \"validate\": invalid syntax")
	})
}
