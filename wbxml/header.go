package wbxml

import (
	"io"
)

type Header struct {
	versionNumber    byte
	publicIdentifier uint32
	charSet          uint32
	stringTable      StringTable
}

func (h *Header) ReadFromBuffer(reader io.ByteReader) error {
	var err error
	h.versionNumber, err = reader.ReadByte()
	if err == nil {
		h.publicIdentifier, err = readMultiByteUint32(reader)
		if err == nil {
			h.charSet, err = readMultiByteUint32(reader)
			if err == nil {
				err = h.stringTable.ReadFromBuffer(reader)
			}
		}
	}

	return err
}

func NewDefaultHeader() Header {
	var header Header
	header.versionNumber = WBXML_1_3
	header.publicIdentifier = uint32(UNKNOWN_PI)
	header.charSet = uint32(CHARSET_UTF8)
	header.stringTable.length = 0
	header.stringTable.content = nil

	return header
}
