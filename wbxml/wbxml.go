package wbxml

import (
	"io"
)

func Decode(reader io.ByteScanner, codeBook *CodeBook) (string, error) {
	decoder := NewDecoder(reader, codeBook)
	err := decoder.header.Read(reader)

	if err == nil {
		return decoder.Decode()
	}

	return "", err
}

func Encode(codeBook *CodeBook, xmlData string, writer io.Writer) error {
	encoder := NewEncoder(codeBook, xmlData, writer)
	err := encoder.Encode()
	return err
}
