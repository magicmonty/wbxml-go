package wbxml

import (
	"bytes"
	"fmt"
	"io"
)

const (
	TAG_BR     byte = 0x05
	TAG_CARD   byte = 0x06
	TAG_XYZ    byte = 0x07
	TAG_DO     byte = 0x08
	TAG_INPUT  byte = 0x09
	TAG_CP2TAG byte = 0x05

	ATTR_STYLE_LIST byte = 0x05
	ATTR_TYPE       byte = 0x06
	ATTR_TYPE_TEXT  byte = 0x07
	ATTR_URL_HTTP   byte = 0x08
	ATTR_NAME       byte = 0x09
	ATTR_KEY        byte = 0x0A

	VALUE_ORG    byte = 0x85
	VALUE_ACCEPT byte = 0x86
)

func makeCodeBook() *CodeBook {
	var codeBook *CodeBook = NewCodeBook()

	var codePage CodePage = NewCodePage("cp", 0)
	codePage.AddTag("BR", TAG_BR)
	codePage.AddTag("CARD", TAG_CARD)
	codePage.AddTag("XYZ", TAG_XYZ)
	codePage.AddTag("DO", TAG_DO)
	codePage.AddTag("INPUT", TAG_INPUT)

	codeBook.AddTagCodePage(codePage)

	codePage = NewCodePage("cp2", 1)
	codePage.AddTag("CP2TAG", TAG_CP2TAG)

	codeBook.AddTagCodePage(codePage)

	codePage = NewCodePage("cp255", 255)
	codeBook.AddTagCodePage(codePage)

	var attributeCodePage AttributeCodePage = NewAttributeCodePage(0)
	attributeCodePage.AddAttribute("STYLE", "LIST", ATTR_STYLE_LIST)
	attributeCodePage.AddAttribute("TYPE", "", ATTR_TYPE)
	attributeCodePage.AddAttribute("TYPE", "TEXT", ATTR_TYPE_TEXT)
	attributeCodePage.AddAttribute("URL", "http://", ATTR_URL_HTTP)
	attributeCodePage.AddAttribute("NAME", "", ATTR_NAME)
	attributeCodePage.AddAttribute("KEY", "", ATTR_KEY)
	attributeCodePage.AddAttributeValue(".org", VALUE_ORG)
	attributeCodePage.AddAttributeValue("ACCEPT", VALUE_ACCEPT)

	codeBook.AddAttributeCodePage(attributeCodePage)

	return codeBook
}

func makeDataBuffer(data ...byte) *bytes.Buffer {
	return bytes.NewBuffer(data)
}

func getDecodeResult(data ...byte) string {
	var result string
	result, _ = Decode(bytes.NewBuffer(data), makeCodeBook())
	return result
}

func printByteStream(r *bytes.Buffer) {
	var (
		result string = ""
	)

	for true {
		b, err := r.ReadByte()
		if err == io.EOF {
			break
		}

		result += fmt.Sprintf("%0.2X ", b)
	}

	fmt.Println(result)
}
