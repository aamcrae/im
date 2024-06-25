package im

import (
	"bytes"
	"fmt"
	"io"
)

const (
	sectionStart = 0xFF
	soiMarker    = 0xD8
	eoiMarker    = 0xD9
	sosMarker    = 0xDA
	dqtMarker    = 0xDB
	dhtMarker    = 0xC4
	sof0Marker   = 0xC0
	sof2Marker   = 0xC2
	app0Marker   = 0xE0
	app1Marker   = 0xE1
	app2Marker   = 0xE2
	app3Marker   = 0xE3
	app4Marker   = 0xE4
	app5Marker   = 0xE5
	app6Marker   = 0xE6
	app7Marker   = 0xE7
	app8Marker   = 0xE8
	app9Marker   = 0xE9
	app10Marker  = 0xEA
	app11Marker  = 0xEB
	app12Marker  = 0xEC
	app13Marker  = 0xED
	app14Marker  = 0xEE
	app15Marker  = 0xEF
	comMarker    = 0xFE
)

type jpegFile struct {
	imageFile
}

func init() {
	registerFileType(jpegFileType)
}

func jpegFileType(f io.ReadSeeker) imageFile {
	b := make([]byte, 2)
	n, err := f.Read(b)
	if err != nil || n != 2 || b[0] != sectionStart || b[1] != soiMarker {
		return nil
	}
	return &jpegFile{}
}

func (f *jpegFile) ReadMeta(im *Imeta, rd io.ReadSeeker) error {
	// Found a JPEG file, read each section until the start of image section is found
	b := make([]byte, 2)
	retError := fmt.Errorf("malformed jpeg")
	for {
		n, err := rd.Read(b)
		if err != nil || n != 2 || b[0] != sectionStart {
			// Malformed jpeg
			return retError
		}
		switch b[1] {
		case eoiMarker, sosMarker:
			// End of image or start of scan, no more metadata.
			return nil
		case soiMarker:
			return retError
		case dqtMarker, sof0Marker, sof2Marker,
			dhtMarker, comMarker:
			if err := skipJpegSection(rd); err != nil {
				return err
			}
		case app0Marker, app1Marker, app2Marker, app3Marker,
			app4Marker, app5Marker, app6Marker, app7Marker,
			app8Marker, app9Marker, app10Marker, app11Marker,
			app12Marker, app13Marker, app14Marker, app15Marker:
			sect, err := readJpegSection(rd)
			if err != nil {
				return err
			}
			if err = jpegApp(im, b[1], sect); err != nil {
				return err
			}
		default:
			fmt.Printf("Unknown section (0x%x)\n", b[1])
			return retError
		}
	}
	return nil
}

func jpegApp(im *Imeta, marker byte, sect []byte) error {
	switch marker {
	case app0Marker:
		// JFIF
	case app1Marker:
		fmt.Printf("App1, tag = %s\n", string(bytes.Split(sect, []byte{0})[0]))
		// Exif or XMP
		if bytes.HasPrefix(sect, []byte("Exif\x00\x00")) {
			if err := im.addExif(sect[6:]); err != nil {
				return err
			}
		} else if bytes.HasPrefix(sect, []byte("http://ns.adobe.com/xap/1.0/\x00")) {
			if err := im.addXmp(sect[29:]); err != nil {
				return err
			}
		}
	case app13Marker:
		// IPTC
		if bytes.HasPrefix(sect, []byte("Photoshop 3.0\x00")) ||
			bytes.HasPrefix(sect, []byte("Photoshop 2.5\x00")) {
			if err := im.addPIR(sect[14:]); err != nil {
				return err
			}
		}
	}
	return nil
}

func skipJpegSection(f io.ReadSeeker) error {
	sz, err := readJpegSectionSize(f)
	if err != nil {
		return err
	}
	f.Seek(int64(sz), io.SeekCurrent)
	return nil
}

func readJpegSection(f io.ReadSeeker) ([]byte, error) {
	sz, err := readJpegSectionSize(f)
	if err != nil {
		return nil, err
	}
	sect := make([]byte, sz)
	_, err = f.Read(sect)
	return sect, err
}

func readJpegSectionSize(f io.ReadSeeker) (int, error) {
	b := make([]byte, 2)
	n, err := f.Read(b)
	if err != nil || n != 2 {
		return -1, err
	}
	sz := (int(b[0]) << 8) + int(b[1])
	// Size includes the 2 byte size value
	if sz < 2 {
		return -1, fmt.Errorf("illegal section size (%d)", sz)
	}
	return sz - 2, nil
}
