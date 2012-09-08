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

func (h *Header) ReadFromBuffer(reader io.ByteReader) {
	h.versionNumber, _ = reader.ReadByte()
	h.publicIdentifier = readMultiByteUint32(reader)
	h.charSet = readMultiByteUint32(reader)
	h.stringTable.ReadFromBuffer(reader)
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
