package im

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

var fileTypes []fileType

// registerFileType registers a callback used to detect a specific file type.
func registerFileType(f fileType) {
	fileTypes = append(fileTypes, f)
}

// ReadFromFile reads from an image file and extracts the metadata
func ReadFromFile(src string) (*Imeta, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Read(f)
}

// ReadFromBuf reads from a buffer and extracts the metadata
func ReadFromBuf(b []byte) (*Imeta, error) {
	return Read(bytes.NewReader(b))
}

// Read reads from a ReadSeeker interface and extracts the metadata
func Read(f io.ReadSeeker) (*Imeta, error) {
	ft := detectFileType(f)
	if ft == nil {
		return nil, fmt.Errorf("unrecognised file type")
	}
	im := &Imeta{}
	if err := ft.ReadMeta(im, f); err != nil {
		return nil, err
	}
	return im, nil
}

// detectFileType walks through the registered image file types until
// one recognises the type of image file.
func detectFileType(f io.ReadSeeker) imageFile {
	// Walk through the decoders until one is successful
	for _, d := range fileTypes {
		f.Seek(0, io.SeekStart)
		ft := d(f)
		if ft != nil {
			return ft
		}
	}
	return nil
}
