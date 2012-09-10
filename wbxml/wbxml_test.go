package wbxml

import (
	"bytes"
	"fmt"
	"testing"
)

func TestDummy(t *testing.T) {

}

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

func MakeCodeBook() *CodeBook {
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

func ExampleEmptyTagWithDifferentNameSpace() {
	fmt.Println(
		getDecodeResult(WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00, SWITCH_PAGE, 0x01, TAG_CP2TAG))
	// OUTPUT: <?xml version="1.0"?>
	// <B:CP2TAG xmlns:B="cp2"/>
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

func ExampleTagWithEmptyTagFromDifferentCodePageAsContent() {
	fmt.Println(
		getDecodeResult(
			WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00,
			TAG_XYZ|TAG_HAS_CONTENT,
			SWITCH_PAGE, 0x01, TAG_CP2TAG,
			END))
	// OUTPUT: <?xml version="1.0"?>
	// <XYZ xmlns="cp" xmlns:B="cp2"><B:CP2TAG/></XYZ>
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

func ExampleTagFromDifferentCodePageWithTextAsContent() {
	fmt.Println(
		getDecodeResult(
			WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00,
			SWITCH_PAGE, 0x01, TAG_CP2TAG|TAG_HAS_CONTENT,
			STR_I, 'X', ' ', '&', ' ', 'Y', 0x00,
			END))
	// OUTPUT: <?xml version="1.0"?>
	// <B:CP2TAG xmlns:B="cp2">X &amp; Y</B:CP2TAG>
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

func ExampleMultipleNestedTagsWithDifferentCodePages() {
	fmt.Println(
		getDecodeResult(
			WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00,
			TAG_XYZ|TAG_HAS_CONTENT,
			SWITCH_PAGE, 0x01, TAG_CP2TAG|TAG_HAS_CONTENT,
			SWITCH_PAGE, 0x00, TAG_DO|TAG_HAS_CONTENT,
			TAG_BR,
			END, END, END))
	// OUTPUT: <?xml version="1.0"?>
	// <XYZ xmlns="cp" xmlns:B="cp2"><B:CP2TAG><DO><BR/></DO></B:CP2TAG></XYZ>
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
	decoder := NewDecoder(
		MakeDataBuffer(
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

// Example from http://www.w3.org/TR/wbxml/#_Toc443384927
func ExampleExtendedWBXMLDecode() {
	fmt.Println(
		getDecodeResult(
			WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8,
			0x12,
			'a', 'b', 'c', 0x00,
			' ', 'E', 'n', 't', 'e', 'r', ' ', 'n', 'a', 'm', 'e', ':', ' ', 0x00,
			TAG_XYZ|TAG_HAS_CONTENT,
			TAG_CARD|TAG_HAS_CONTENT|TAG_HAS_ATTRIBUTES,
			ATTR_NAME, STR_T, 0x00, ATTR_STYLE_LIST, END,
			TAG_DO|TAG_HAS_ATTRIBUTES,
			ATTR_TYPE, VALUE_ACCEPT,
			ATTR_URL_HTTP, STR_I, 'x', 'y', 'z', 0x00, VALUE_ORG, STR_I, '/', 's', 0x00, END,
			STR_T, 0x04,
			TAG_INPUT|TAG_HAS_ATTRIBUTES,
			ATTR_TYPE_TEXT, ATTR_KEY, STR_I, 'N', 0x00, END,
			END,
			END))
	// OUTPUT: <?xml version="1.0"?>
	// <XYZ><CARD NAME="abc" STYLE="LIST"><DO TYPE="ACCEPT" URL="http://xyz.org/s"/> Enter name: <INPUT TYPE="TEXT" KEY="N"/></CARD></XYZ>
}

func ExampleGetNameSpaceDeclarations() {
	decoder := NewDecoder(MakeDataBuffer(0x00), MakeCodeBook())

	decoder.usedNamespaces[0] = true
	decoder.usedNamespaces[1] = true
	decoder.usedNamespaces[255] = true
	fmt.Println(decoder.getNameSpaceDeclarations())
	// OUTPUT:  xmlns="cp" xmlns:B="cp2" xmlns:IV="cp255"
}

func TestGetNameSpaceDeclarationsShouldReturnEmptyStringIfOnlyCP0IsSelected(t *testing.T) {
	decoder := NewDecoder(MakeDataBuffer(0x00), MakeCodeBook())

	decoder.usedNamespaces[0] = true
	if decoder.getNameSpaceDeclarations() != "" {
		t.Errorf("NameSpace declaration should be emty but was \"%s\"", decoder.getNameSpaceDeclarations())
	}
}

func TestGetNameSpaceDeclarationsShouldReturnEmptyStringINoCPIsActive(t *testing.T) {
	decoder := NewDecoder(MakeDataBuffer(0x00), MakeCodeBook())

	if decoder.getNameSpaceDeclarations() != "" {
		t.Errorf("NameSpace declaration should be emty but was \"%s\"", decoder.getNameSpaceDeclarations())
	}
}
