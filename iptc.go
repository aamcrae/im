package im

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type iptcRec struct {
	Tag byte
	Rec byte
	DS  byte
	Len uint16
}

// addIptc parses the blob and extracts the IPTC data
func (im *Imeta) addIptc(b []byte) error {
	rd := bytes.NewReader(b)
	for rd.Len() > 0 {
		var rec iptcRec
		err := binary.Read(rd, binary.BigEndian, &rec)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if rec.Tag != 0x1C {
			return fmt.Errorf("unexpected tag in IPTC record (0x%02x)", rec.Tag)
		}
		var l int
		var data []byte
		// If the length has the top bit set, this indicates an
		// extended DataSet tag. The bottom 15 bits of the length is used as the length
		// of a following field that defines the length of the data block.
		if (rec.Len & 0x8000) != 0 {
			lenRd := int(rec.Len & 0x7FFF)
			lenData := make([]byte, lenRd)
			err := binary.Read(rd, binary.BigEndian, &lenData)
			if err != nil {
				return err
			}
			// Calculate the length from the extended data length field
			for _, b := range lenData {
				l = (l << 8) + int(b)
				// Sanity check.
				if l > rd.Len() {
					return fmt.Errorf("bad length data in extended dataset (rec %d, DS %d)", rec.Rec, rec.DS)
				}
			}
		} else {
			l = int(rec.Len)
		}
		if l != 0 {
			data = make([]byte, l)
			err := binary.Read(rd, binary.BigEndian, &data)
			if err != nil {
				return err
			}
		}
		fmt.Printf("IPTC rec %d, DS %d, data = %v\n", rec.Rec, rec.DS, data)
	}
	return nil
}
