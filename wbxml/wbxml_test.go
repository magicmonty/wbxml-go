package wbxml

import (
	"bytes"
	"fmt"
	"testing"
)

func TestDummy(t *testing.T) {

}

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

	codeBook.AddTagCodePage(codePage)

	return codeBook
}

func MakeDataBuffer(data ...byte) *bytes.Buffer {
	return bytes.NewBuffer(data)
}

func getDecodeResult(data ...byte) string {
	var result string
	result, _ = Decode(bytes.NewBuffer(data), MakeCodeBook())
	return result
}

func ExampleEmptyTag() {
	fmt.Println(
		getDecodeResult(WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00, TAG_XYZ))
	// OUTPUT: <?xml version="1.0"?>
	// <XYZ/>
}

func ExampleEmptyLiteralTag() {
	fmt.Println(
		getDecodeResult(
			WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8,
			0x04, 'X', 'Y', 'Z', 0x00,
			LITERAL, 0x00))
	// OUTPUT: <?xml version="1.0"?>
	// <XYZ/>
}

func ExampleTagWithEmptyTagAsContent() {
	fmt.Println(
		getDecodeResult(
			WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00,
			TAG_XYZ|TAG_HAS_CONTENT, TAG_CARD, END))
	// OUTPUT: <?xml version="1.0"?>
	// <XYZ><CARD/></XYZ>
}

func ExampleTagWithTextAsContent() {
	fmt.Println(
		getDecodeResult(
			WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00,
			TAG_XYZ|TAG_HAS_CONTENT,
			STR_I, 'X', ' ', '&', ' ', 'Y', 0x00,
			END))
	// OUTPUT: <?xml version="1.0"?>
	// <XYZ>X &amp; Y</XYZ>
}

func ExampleMultipleNestedTags() {
	fmt.Println(
		getDecodeResult(
			WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00,
			TAG_XYZ|TAG_HAS_CONTENT, TAG_CARD|TAG_HAS_CONTENT, TAG_DO|TAG_HAS_CONTENT,
			TAG_BR,
			END, END, END))
	// OUTPUT: <?xml version="1.0"?>
	// <XYZ><CARD><DO><BR/></DO></CARD></XYZ>
}

func ExampleReadMultiByteUint32() {
	var (
		result uint32
	)

	result, _ = readMultiByteUint32(MakeDataBuffer(0x81, 0x20))
	fmt.Printf("%d\n", result)
	result, _ = readMultiByteUint32(MakeDataBuffer(0x60))
	fmt.Printf("%d\n", result)
	// OUTPUT: 160
	// 96
}

func ExampleDecodeEntity() {
	decoder, _ := NewDecoder(
		MakeDataBuffer(
			WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00,
			ENTITY, 0x81, 0x20,
			ENTITY, 0x60),
		MakeCodeBook())
	var result string
	result, _ = decoder.decodeEntity()
	fmt.Println(result)
	result, _ = decoder.decodeEntity()
	fmt.Println(result)
	// OUTPUT: &#160;
	// &#96;
}

// Example from http://www.w3.org/TR/wbxml/#_Toc443384926
func ExampleSimpleWBXMLDecode() {
	fmt.Println(
		getDecodeResult(
			WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00,
			TAG_XYZ|TAG_HAS_CONTENT, TAG_CARD|TAG_HAS_CONTENT,
			STR_I, ' ', 'X', ' ', '&', ' ', 'Y', 0x00,
			TAG_BR,
			STR_I, ' ', 'X', 0x00,
			ENTITY, 0x81, 0x20,
			STR_I, '=', 0x00,
			ENTITY, 0x81, 0x20,
			STR_I, '1', ' ', 0x00,
			END, END))
	// OUTPUT: <?xml version="1.0"?>
	// <XYZ><CARD> X &amp; Y<BR/> X&#160;=&#160;1 </CARD></XYZ>
}
