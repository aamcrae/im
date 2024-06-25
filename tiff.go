package im

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type tiffFile struct {
	imageFile
	order binary.ByteOrder
}

type tiffTag struct {
	Id      uint16
	DType   uint16
	DCount  uint32
	Doffset [4]byte
}

func init() {
	registerFileType(tiffFileType)
}

func tiffFileType(f io.ReadSeeker) imageFile {
	b := make([]byte, 2)
	n, err := f.Read(b)
	if err != nil || n != 2 {
		return nil
	}
	var order binary.ByteOrder
	if bytes.Equal(b, []byte{'I', 'I'}) {
		order = binary.LittleEndian
	} else if bytes.Equal(b, []byte{'M', 'I'}) {
		order = binary.LittleEndian
	} else {
		// Not a TIFF file
		return nil
	}
	// Read magic number of 42
	var magic uint16
	err = binary.Read(f, order, &magic)
	if err != nil || magic != 42 {
		return nil
	}
	return &tiffFile{order: order}
}

func (t *tiffFile) ReadMeta(im *Imeta, f io.ReadSeeker) error {
	// Valid TIFF file, read IFDs
	tifds, err := readOffsetTiffTags(f, t.order)
	if err != nil {
		return err
	}
	for _, ifds := range tifds {
		fmt.Printf("IFD, length %d\n", len(ifds))
		for _, tag := range ifds {
			fmt.Printf("Tag 0x%04x, type %d, count %d\n", tag.Id, tag.DType, tag.DCount)
		}
	}
	return nil
}

func readOffsetTiffTags(rd io.ReadSeeker, order binary.ByteOrder) ([][]tiffTag, error) {
	// Read offset of IFD first, then seek to IFDs
	var offs uint32
	if err := binary.Read(rd, order, &offs); err != nil {
		return nil, err
	}
	return readTiffTags(rd, order, int64(offs))
}

func readTiffTags(rd io.ReadSeeker, order binary.ByteOrder, offset int64) ([][]tiffTag, error) {
	var ifds [][]tiffTag
	for offset != 0 {
		// Seek to next IFD
		if _, err := rd.Seek(offset, io.SeekStart); err != nil {
			return nil, err
		}
		var ndirs uint16
		var tags []tiffTag
		if err := binary.Read(rd, order, &ndirs); err != nil {
			return nil, err
		}
		fmt.Printf("IFD, ndirs = %d, offs = %d\n", ndirs, offset)
		for ; ndirs > 0; ndirs-- {
			var tag tiffTag
			if err := binary.Read(rd, order, &tag); err != nil {
				return nil, err
			}
			tags = append(tags, tag)
		}
		ifds = append(ifds, tags)
		var newOffset uint32
		if err := binary.Read(rd, order, &newOffset); err != nil {
			return nil, err
		}
		offset = int64(newOffset)
	}
	return ifds, nil
}
