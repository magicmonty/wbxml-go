package wbxml

import (
	"io"
)

func decodeTagWithContentAndAttributes(reader io.ByteScanner, codeBook *CodeBook, currentCodePage CodePage) string {
	return ""
}

func decodeTagWithContent(reader io.ByteScanner, codeBook *CodeBook, currentCodePage CodePage) string {
	return ""
}

func decodeEmptyTagWithAttributes(reader io.ByteScanner, codeBook *CodeBook, currentCodePage CodePage) string {
	return ""
}

func decodeEmptyTag(reader io.ByteScanner, codeBook *CodeBook, currentCodePage CodePage) string {
	var (
		nextByte byte
	)

	nextByte, _ = reader.ReadByte()

	if currentCodePage.HasTagCode(nextByte) {
		return "<" + currentCodePage.Tags[nextByte] + "/>"
	}

	return ""
}

func decodeElement(reader io.ByteScanner, codeBook *CodeBook, currentCodePage CodePage) string {
	var (
		nextByte byte
	)

	nextByte, _ = reader.ReadByte()
	reader.UnreadByte()

	if nextByte&EXT_I_0 != 0 {
		if nextByte&EXT_0 != 0 {
			return decodeTagWithContentAndAttributes(reader, codeBook, currentCodePage)
		} else {
			return decodeTagWithContent(reader, codeBook, currentCodePage)
		}
	} else if nextByte&EXT_0 != 0 {
		return decodeEmptyTagWithAttributes(reader, codeBook, currentCodePage)
	} else {
		return decodeEmptyTag(reader, codeBook, currentCodePage)
	}

	return ""
}

func decodeBody(reader io.ByteScanner, codeBook *CodeBook) string {
	var currentCodePage CodePage = codeBook.CodePages[0]

	return decodeElement(reader, codeBook, currentCodePage)
}

func Decode(reader io.ByteScanner, codeBook *CodeBook) string {
	var (
		result string = "<?xml version=\"1.0\"?>\n"
		header Header
	)

	header.ReadFromBuffer(reader)
	result += decodeBody(reader, codeBook)

	return result
}
