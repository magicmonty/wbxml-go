package wbxml

import (
	"io"
)

func Decode(reader io.ByteScanner, codeBook *CodeBook) (string, error) {
	decoder := NewDecoder(reader, codeBook)
	err := decoder.header.ReadFromBuffer(reader)

	if err == nil {
		return decoder.decodeBody()
	}

	return "", err
}
