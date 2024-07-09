package im

import (
	"io"
)

type Type int

const (
	TypeNone Type = iota
	TypeByte
	TypeAscii
	TypeShort    // 16 bit
	TypeLong     // 32 bit
	TypeRational // 2 x 32 bit unsigned
	TypeSByte
	TypeUndefine
	TypeSShort
	TypeSLong
	TypeSRational // 2 x 32 bit signed
	TypeFloat
	TypeDouble
	TypeXmp
)

type Rational struct {
	Num, Denom uint32
}

type SRational struct {
	Num, Denom int32
}

type Value struct {
	Type  Type
	Count int // Count of items
	v     any
}

type ValueError struct {
	Err  string
	Type Type
}

type Item struct {
	val Value
}

type imageFile interface {
	ReadMeta(*Imeta, io.ReadSeeker) error
}

type fileType func(io.ReadSeeker) imageFile

type group struct {
	items map[string]*Item
}

type Imeta struct {
	groups map[string]*group
}
