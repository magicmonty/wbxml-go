package wbxml

import (
	"io"
)

func decodeTagWithContentAndAttributes(reader io.ByteScanner, codeBook *CodeBook, currentCodePage CodePage) string {
	return ""
}

func decodeTagWithContent(reader io.ByteScanner, codeBook *CodeBook, currentCodePage CodePage) string {
	var (
		result     string = ""
		nextByte   byte
		tagCode    byte
		currentTag string
	)

	nextByte, _ = reader.ReadByte()
	tagCode = nextByte &^ TAG_HAS_CONTENT
	if currentCodePage.HasTagCode(tagCode) {
		currentTag = currentCodePage.Tags[tagCode]
		result = "<" + currentTag + ">"
		result += decodeElement(reader, codeBook, currentCodePage)
		nextByte, _ = reader.ReadByte()
		if nextByte == END {
			result += "</" + currentTag + ">"
		}
	}

	return result
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

	if nextByte&TAG_HAS_CONTENT != 0 {

		if nextByte&TAG_HAS_ATTRIBUTES != 0 {
			return decodeTagWithContentAndAttributes(reader, codeBook, currentCodePage)
		} else {
			return decodeTagWithContent(reader, codeBook, currentCodePage)
		}
	} else if nextByte&TAG_HAS_ATTRIBUTES != 0 {
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
