package wbxml

import (
	"io"
)

type Header struct {
	versionNumber    byte
	publicIdentifier byte
	charSet          byte
	stringTable      StringTable
}

func (h *Header) ReadFromBuffer(reader io.ByteReader) {
	h.versionNumber, _ = reader.ReadByte()
	h.publicIdentifier, _ = reader.ReadByte()
	h.charSet, _ = reader.ReadByte()
	h.stringTable.ReadFromBuffer(reader)
}

func NewDefaultHeader() Header {
	var header Header
	header.versionNumber = WBXML_1_3
	header.publicIdentifier = UNKNOWN_PI
	header.charSet = CHARSET_UTF8
	header.stringTable.length = 0
	header.stringTable.content = nil

	return header
}
