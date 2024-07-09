package im

import (
	"bytes"
	"encoding/binary"
	"io"
)

// String returns the name of the type
func (t Type) String() string {
	switch t {
	case TypeNone:
		return "<nil>"
	case TypeByte:
		return "Byte"
	case TypeAscii:
		return "Ascii"
	case TypeShort:
		return "Short"
	case TypeLong:
		return "Long"
	case TypeRational:
		return "Rational"
	case TypeSByte:
		return "SByte"
	case TypeUndefine:
		return "Undefine"
	case TypeSShort:
		return "SShort"
	case TypeSLong:
		return "SLong"
	case TypeSRational:
		return "SRational"
	case TypeFloat:
		return "Float"
	case TypeDouble:
		return "Double"
	case TypeXmp:
		return "XmpText"
	default:
		panic("unknown type")
	}
}

// Size returns the size in bytes for this type.
func (t Type) Size() int {
	switch t {
	case TypeNone:
		return 0
	case TypeByte, TypeAscii, TypeSByte, TypeUndefine, TypeXmp:
		return 1
	case TypeShort, TypeSShort:
		return 2
	case TypeLong, TypeSLong, TypeFloat:
		return 4
	case TypeRational, TypeSRational, TypeDouble:
		return 8
	default:
		panic("unknown type")
	}
}

// read reads 1 or more items of this type from the reader.
func (t Type) read(count int, order binary.ByteOrder, rd io.Reader) (any, error) {
	switch t {
	case TypeNone:
		return nil, nil
	case TypeByte:
		return readData[byte](count, order, rd)
	case TypeAscii:
		return readString(count, order, rd)
	case TypeShort:
		return readData[uint16](count, order, rd)
	case TypeLong:
		return readData[uint32](count, order, rd)
	case TypeRational:
		return readRational[Rational](count, order, rd)
	case TypeSByte:
		return readData[int8](count, order, rd)
	case TypeUndefine:
		return readData[byte](count, order, rd)
	case TypeSShort:
		return readData[int16](count, order, rd)
	case TypeSLong:
		return readData[int32](count, order, rd)
	case TypeSRational:
		return readRational[SRational](count, order, rd)
	case TypeFloat:
		return readData[float32](count, order, rd)
	case TypeDouble:
		return readData[float64](count, order, rd)
	case TypeXmp:
		return readString(count, order, rd)
	default:
		panic("unknown type")
	}
}

// IsExif returns true if this type is a valid EXIF type
func (t Type) IsExif() bool {
	switch t {
	case TypeByte, TypeAscii, TypeSByte, TypeUndefine,
		TypeShort, TypeSShort, TypeLong, TypeSLong,
		TypeRational, TypeSRational, TypeFloat:
		return true
	default:
		return false
	}
}

// IsXmp returns true if the type is an XMP string type
func (t Type) IsXmp() bool {
	return t == TypeXmp
}

// readString reads a string of the size indicated. The string may include
// trailing nulls, which are stripped (and will be re-added if this value is re-written).
func readString(count int, order binary.ByteOrder, rd io.Reader) (string, error) {
	if v, err := readData[byte](count, order, rd); err != nil {
		return "", err
	} else {
		// Convert to byte slice if necessary
		var b []byte
		if count == 1 {
			b = []byte{v.(byte)}
		} else {
			b = v.([]byte)
		}
		// Trip trailing nulls from the string
		return string(bytes.TrimRight(b, "\000")), nil
	}
}

// readData reads a list of one or more of these types
func readData[T byte | int8 | int16 | int32 | uint16 | uint32 | float32 | float64](count int, order binary.ByteOrder, rd io.Reader) (any, error) {
	entry := make([]T, count)
	if err := binary.Read(rd, order, &entry); err != nil {
		return nil, err
	}
	if count == 1 {
		return entry[0], nil
	} else {
		return entry, nil
	}
}

// readRational reads a list of one or more rational types (2 32 bit values used as
// fractional numerator and denominator).
func readRational[T Rational | SRational](count int, order binary.ByteOrder, rd io.Reader) (any, error) {
	var sl []T
	for i := 0; i < count; i++ {
		var entry T
		err := binary.Read(rd, order, &entry)
		if err != nil {
			return sl, err
		}
		sl = append(sl, entry)
	}
	if count == 1 {
		return sl[0], nil
	} else {
		return sl, nil
	}
}
