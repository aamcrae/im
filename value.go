package im

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

func (v Value) Value() any {
	switch v.Type {
	case TypeNone:
		return nil
	case TypeByte:
		return v.v
	case TypeXmp:
		return v.v
	default:
		panic("unknown type")
	}
}

func (v Value) String() string {
	return fmt.Sprintf("%v", v.v)
}

func (v Value) Size() int {
	return v.Count * v.Type.Size()
}

func NewValue(t Type, v any) (Value, error) {
	val := Value{Type: t}
	switch t {
	case TypeNone:
	case TypeByte:
	case TypeAscii:
		// Check string is null terminated
		s := v.(string)
		if !strings.HasSuffix(s, "\000") {
			s = s + "\000"
		}
		val.Count = len(s)
		val.v = s
	case TypeXmp:
		s := v.(string)
		val.Count = 1
		val.v = s
	default:
		panic("unknown type")
	}
	return val, nil
}

func readValue(t Type, count int, order binary.ByteOrder, rd io.Reader) (Value, error) {
	val := Value{Type: t, Count: count}
	var err error
	val.v, err = t.read(count, order, rd)
	return val, err
}

func (e *ValueError) Error() string {
	return "imeta: " + e.Err + " on " + e.Type.String()
}
