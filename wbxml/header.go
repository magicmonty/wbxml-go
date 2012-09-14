package wbxml

import (
	"io"
)

type Header struct {
	versionNumber    byte
	publicIdentifier uint32
	charSet          uint32
	charSetAsString  string
	stringTable      *StringTable
}

func (h *Header) Read(reader io.ByteReader) error {
	var err error
	h.versionNumber, err = reader.ReadByte()
	if err == nil {
		h.publicIdentifier, err = readMultiByteUint32(reader)
		if err == nil {
			h.charSet, err = readMultiByteUint32(reader)
			if err == nil {
				h.charSetAsString, _ = GetCharsetStringByCode(h.charSet)
				err = h.stringTable.Read(reader)
			}
		}
	}

	return err
}

func (h *Header) Write(writer io.Writer) error {
	_, err := writer.Write([]byte{h.versionNumber})
	if err == nil {
		err = writeMultiByteUint32(writer, h.publicIdentifier)
		if err == nil {
			err = writeMultiByteUint32(writer, h.charSet)
			if err == nil {
				err = h.stringTable.Write(writer)
			}
		}
	}
	return err
}

func NewDefaultHeader() *Header {
	header := new(Header)
	header.versionNumber = WBXML_1_3
	header.publicIdentifier = uint32(UNKNOWN_PI)
	header.charSet = uint32(CHARSET_UTF8)
	header.stringTable = NewStringTable()
	header.charSetAsString = ""

	return header
}
