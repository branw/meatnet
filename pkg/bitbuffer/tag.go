package bitbuffer

import (
	"errors"
	"reflect"
	"strconv"
)

const (
	// TagBitWidth allows the bit-width of primitive struct fields to be overridden
	TagBitWidth = "bbwidth"
	// TagValidate allows the value to be validated directly after parsing
	TagValidate = "bbvalidate"
)

type ParsedTag struct {
	BitWidth uint64
	Validate bool
}

func ParseTag(tag reflect.StructTag) (ParsedTag, error) {
	bitWidth, err := ParseBitWidthTag(tag)
	if err != nil {
		return ParsedTag{}, err
	}

	validate, err := ParseValidateTag(tag)
	if err != nil {
		return ParsedTag{}, err
	}

	tags := ParsedTag{
		BitWidth: bitWidth,
		Validate: validate,
	}
	return tags, nil
}

func (tags ParsedTag) BitWidthOrDefault(defaultBitWidth uint64) uint64 {
	if tags.BitWidth > 0 {
		return tags.BitWidth
	}
	return defaultBitWidth
}

func ParseBitWidthTag(tag reflect.StructTag) (uint64, error) {
	value, present := tag.Lookup(TagBitWidth)
	if !present {
		return 0, nil
	}

	width, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return 0, err
	}
	// We use the zero value to indicate no width was specified, so we should
	// warn the user if they pass it explicitly
	if width == 0 {
		return 0, errors.New("bit width must be greater than zero")
	}

	return width, nil
}

func ParseValidateTag(tag reflect.StructTag) (bool, error) {
	value, present := tag.Lookup(TagValidate)
	if !present {
		return false, nil
	}

	validate, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}
	return validate, nil
}
