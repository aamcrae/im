package im

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// addPIR adds a Photoshop Image Resource blob, which
// may include IPTC records, thumbnails or other metadata.
func (im *Imeta) addPIR(b []byte) error {
	if len(b) >= 12 && bytes.HasPrefix(b, []byte("8BIM\004")) {
		id := b[5]
		// Next field is 'pascal' string, first
		// byte is length, with padding to make field even length
		lb := int(b[6]) + 1
		lb += lb % 2
		// Check to ensure the block is large enough to read the data block length
		if len(b) < lb+10 {
			return fmt.Errorf("short PIR header")
		}
		b = b[lb+6:] // Move past header to data block length
		sz := binary.BigEndian.Uint32(b)
		b = b[4:]
		if int(sz) != len(b) {
			return fmt.Errorf("size mismatch on PIR data")
		}
		switch id {
		case 0x04:
			return im.addIptc(b)
		case 0x0C:
			fmt.Printf("Thumbnail block\n")
			return nil
		case 0x0F:
			fmt.Printf("ICC Profile block\n")
			return nil
		case 0x22:
			fmt.Printf("EXIF block\n")
			return nil
		case 0x24:
			fmt.Printf("XMP block\n")
			return nil
		}
	}
	return fmt.Errorf("unknown PIR header")
}
