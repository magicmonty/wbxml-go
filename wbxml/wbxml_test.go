package wbxml

import (
	"bytes"
	"fmt"
	"testing"
)

func TestDummy(t *testing.T) {

}

/*
func _ExampleWBXMLDecode() {
	var data []byte = make([]byte, 34)
	data = []byte{
		0x01, 0x01, 0x03, 0x00, 0x47, 0x46, 0x03, 0x20,
		0x58, 0x20, 0x26, 0x20, 0x59, 0x00, 0x05, 0x03,
		0x20, 0x58, 0x00, 0x02, 0x81, 0x20, 0x03, 0x3D,
		0x00, 0x02, 0x81, 0x20, 0x03, 0x31, 0x20, 0x00,
		0x01, 0x01}

	var b *bytes.Buffer = bytes.NewBuffer(data)

	var codeBook *CodeBook
	*codeBook = NewCodeBook()
	var codePage CodePage = NewCodePage("", 0)

	codePage.AddTag("BR", 0x05)
	codePage.AddTag("CARD", 0x06)
	codePage.AddTag("XYZ", 0x07)

	fmt.Println(Decode(b, codeBook))
	// OUTPUT: <?xml version="1.0"?>
	// <XYZ><CARD> X &amp; Y<BR/>X&nbsp;=&nbsp;1</CARD></XYZ>
}
*/

const (
	TAG_BR   byte = 0x05
	TAG_CARD byte = 0x06
	TAG_XYZ  byte = 0x07
	TAG_DO   byte = 0x08
)

func MakeCodeBook() *CodeBook {
	var codeBook *CodeBook = NewCodeBook()

	var codePage CodePage = NewCodePage("cp", 0)
	codePage.AddTag("BR", TAG_BR)
	codePage.AddTag("CARD", TAG_CARD)
	codePage.AddTag("XYZ", TAG_XYZ)
	codePage.AddTag("DO", TAG_DO)

	codeBook.AddCodePage(codePage)

	return codeBook
}

func MakeDataBuffer(data ...byte) *bytes.Buffer {
	return bytes.NewBuffer(data)
}

func ExampleEmptyTag() {
	fmt.Println(
		Decode(
			MakeDataBuffer(WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00, TAG_XYZ),
			MakeCodeBook()))
	// OUTPUT: <?xml version="1.0"?>
	// <XYZ/>
}

func ExampleTagWithEmptyTagAsContent() {
	fmt.Println(
		Decode(
			MakeDataBuffer(WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00, TAG_XYZ|EXT_I_0, TAG_CARD, END),
			MakeCodeBook()))
	// OUTPUT: <?xml version="1.0"?>
	// <XYZ><CARD/></XYZ>
}

func ExampleMultipleNestedTags() {
	fmt.Println(
		Decode(
			MakeDataBuffer(
				WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00,
				TAG_XYZ|EXT_I_0, TAG_CARD|EXT_I_0, TAG_DO|EXT_I_0, TAG_BR, END, END, END),
			MakeCodeBook()))
	// OUTPUT: <?xml version="1.0"?>
	// <XYZ><CARD><DO><BR/></DO></CARD></XYZ>
}
