package wbxml

import (
	"bytes"
)

func decodeTagWithContentAndAttributes(b *bytes.Buffer, codeBook *CodeBook, currentCodePage CodePage) string {
	return ""
}

func decodeTagWithContent(b *bytes.Buffer, codeBook *CodeBook, currentCodePage CodePage) string {
	return ""
}

func decodeEmptyTagWithAttributes(b *bytes.Buffer, codeBook *CodeBook, currentCodePage CodePage) string {
	return ""
}

func decodeEmptyTag(b *bytes.Buffer, codeBook *CodeBook, currentCodePage CodePage) string {
	var nextByte byte
	nextByte, _ = b.ReadByte()

	if currentCodePage.HasCode(nextByte) {
		return "<" + currentCodePage.Codes[nextByte] + "/>"
	}
	return ""
}

func decodeElement(b *bytes.Buffer, codeBook *CodeBook, currentCodePage CodePage) string {
	var (
		nextByte byte
	)

	nextByte, _ = b.ReadByte()
	b.UnreadByte()

	if nextByte&EXT_I_0 != 0 {
		if nextByte&EXT_0 != 0 {
			return decodeTagWithContentAndAttributes(b, codeBook, currentCodePage)
		} else {
			return decodeTagWithContent(b, codeBook, currentCodePage)
		}
	} else if nextByte&EXT_0 != 0 {
		return decodeEmptyTagWithAttributes(b, codeBook, currentCodePage)
	} else {
		return decodeEmptyTag(b, codeBook, currentCodePage)
	}

	return ""
}

func decodeBody(b *bytes.Buffer, codeBook *CodeBook) string {
	var currentCodePage CodePage = codeBook.CodePages[0]

	return decodeElement(b, codeBook, currentCodePage)
}

func Decode(b *bytes.Buffer, codeBook *CodeBook) string {
	var (
		result string = "<?xml version=\"1.0\"?>\n"
		header Header
	)

	header.ReadFromBuffer(b)
	result += decodeBody(b, codeBook)

	return result
}
