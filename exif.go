package im

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// addExif reads the EXIF tags from the blob and saves the tag entries.
func (im *Imeta) addExif(b []byte) error {
	rd := bytes.NewReader(b)
	var align uint16
	if err := binary.Read(rd, binary.BigEndian, &align); err != nil {
		return err
	}
	var order binary.ByteOrder
	if align == 0x4949 {
		order = binary.LittleEndian
	} else if align == 0x4D4D {
		order = binary.BigEndian
	} else {
		return fmt.Errorf("unknown TIFF byte align")
	}
	var version uint16
	if err := binary.Read(rd, order, &version); err != nil {
		return err
	}
	if version != 0x2A {
		return fmt.Errorf("unexpected TIFF version (0x%x)", version)
	}
	tifds, err := readOffsetTiffTags(rd, order)
	if err != nil {
		return err
	}
	if err := im.exifScanIFDs(b, tifds, order, "Image"); err != nil {
		return err
	}
	return nil
}

// exifScanIFDs reads the Image File Directory (IFD) entries
func (im *Imeta) exifScanIFDs(b []byte, tifds [][]tiffTag, order binary.ByteOrder, name string) error {
	for i, tags := range tifds {
		fmt.Printf("%d tags in IFD %d\n", len(tags), i+1)
		if err := im.exifIFD(b, tags, order, name); err != nil {
			return err
		}
	}
	return nil
}

// exifIFD processes a single IFD containing a list of tags
func (im *Imeta) exifIFD(b []byte, tags []tiffTag, order binary.ByteOrder, name string) error {
	for _, tag := range tags {
		val, err := tagToValue(b, order, &tag)
		if err != nil {
			return err
		}
		fmt.Printf("%s tag 0x%04x, type %s, val = %s\n", name, tag.Id, val.Type.String(), val.String())
		if tag.Id == 0x8769 {
			ifdRd := bytes.NewReader(b)
			offs := order.Uint32(tag.Doffset[:])
			fmt.Printf("SubIFD, offset %d\n", offs)
			tifds, err := readTiffTags(ifdRd, order, int64(offs))
			if err != nil {
				return err
			}
			im.exifScanIFDs(b, tifds, order, "Photo")
		}
	}
	return nil
}

// tagToValue converts an EXIF tag value to a Value
func tagToValue(b []byte, order binary.ByteOrder, t *tiffTag) (Value, error) {
	vt := Type(t.DType)
	if !vt.IsExif() {
		return Value{}, fmt.Errorf("unknown datatype (tag 0x%04x)", t.Id)
	}
	if vt.Size()*int(t.DCount) <= len(t.Doffset) {
		// The offset field contains the data
		return readValue(vt, int(t.DCount), order, bytes.NewReader(t.Doffset[:]))
	} else {
		// The offset field points to a separate location
		offs := order.Uint32(t.Doffset[:])
		rd := bytes.NewReader(b)
		rd.Seek(int64(offs), io.SeekStart)
		return readValue(vt, int(t.DCount), order, rd)
	}
}
