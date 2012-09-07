package wbxml

import (
	"bytes"
)

type Header struct {
	versionNumber    byte
	publicIdentifier byte
	charSet          byte
	stringTable      StringTable
}

func (h *Header) ReadFromBuffer(b *bytes.Buffer) {
	h.versionNumber, _ = b.ReadByte()
	h.publicIdentifier, _ = b.ReadByte()
	h.charSet, _ = b.ReadByte()
	h.stringTable.ReadFromBuffer(b)
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
