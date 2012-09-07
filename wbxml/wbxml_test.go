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

	codePage.AddItem("BR", 0x05)
	codePage.AddItem("CARD", 0x06)
	codePage.AddItem("XYZ", 0x07)

	fmt.Println(Decode(b, codeBook))
	// OUTPUT: <?xml version="1.0"?>
	// <XYZ><CARD> X &amp; Y<BR/>X&nbsp;=&nbsp;1</CARD></XYZ>
}
*/

func ExampleEmptyTag() {
	var data []byte = make([]byte, 34)
	data = []byte{0x01, 0x01, 0x03, 0x00, 0x07, 0x00, 0x01, 0x01}

	var b *bytes.Buffer = bytes.NewBuffer(data)

	var codeBook *CodeBook = NewCodeBook()
	var codePage CodePage = NewCodePage("cp", 0)

	codePage.AddItem("XYZ", 0x07)
	codeBook.AddCodePage(codePage)

	fmt.Println(Decode(b, codeBook))
	// OUTPUT: <?xml version="1.0"?>
	// <XYZ/>
}
