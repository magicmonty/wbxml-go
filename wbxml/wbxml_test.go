package wbxml

import (
	"fmt"
)

func ExampleEmptyTag() {
	fmt.Println(
		getDecodeResult(WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00, TAG_XYZ))
	// OUTPUT: <?xml version="1.0" encoding="utf-8"?>
	// <XYZ/>
}

func ExampleEmptyTagWithDifferentNameSpace() {
	fmt.Println(
		getDecodeResult(WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00, SWITCH_PAGE, 0x01, TAG_CP2TAG))
	// OUTPUT: <?xml version="1.0" encoding="utf-8"?>
	// <B:CP2TAG xmlns:B="cp2"/>
}

func ExampleEmptyLiteralTag() {
	fmt.Println(
		getDecodeResult(
			WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8,
			0x04, 'X', 'Y', 'Z', 0x00,
			LITERAL, 0x00))
	// OUTPUT: <?xml version="1.0" encoding="utf-8"?>
	// <XYZ/>
}

func ExampleTagWithEmptyTagAsContent() {
	fmt.Println(
		getDecodeResult(
			WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00,
			TAG_XYZ|TAG_HAS_CONTENT, TAG_CARD, END))
	// OUTPUT: <?xml version="1.0" encoding="utf-8"?>
	// <XYZ><CARD/></XYZ>
}

func ExampleTagWithEmptyTagFromDifferentCodePageAsContent() {
	fmt.Println(
		getDecodeResult(
			WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00,
			TAG_XYZ|TAG_HAS_CONTENT,
			SWITCH_PAGE, 0x01, TAG_CP2TAG,
			END))
	// OUTPUT: <?xml version="1.0" encoding="utf-8"?>
	// <XYZ xmlns="cp" xmlns:B="cp2"><B:CP2TAG/></XYZ>
}

func ExampleTagWithTextAsContent() {
	fmt.Println(
		getDecodeResult(
			WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00,
			TAG_XYZ|TAG_HAS_CONTENT,
			STR_I, 'X', ' ', '&', ' ', 'Y', 0x00,
			END))
	// OUTPUT: <?xml version="1.0" encoding="utf-8"?>
	// <XYZ>X &amp; Y</XYZ>
}

func ExampleTagFromDifferentCodePageWithTextAsContent() {
	fmt.Println(
		getDecodeResult(
			WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00,
			SWITCH_PAGE, 0x01, TAG_CP2TAG|TAG_HAS_CONTENT,
			STR_I, 'X', ' ', '&', ' ', 'Y', 0x00,
			END))
	// OUTPUT: <?xml version="1.0" encoding="utf-8"?>
	// <B:CP2TAG xmlns:B="cp2">X &amp; Y</B:CP2TAG>
}

func ExampleMultipleNestedTags() {
	fmt.Println(
		getDecodeResult(
			WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00,
			TAG_XYZ|TAG_HAS_CONTENT, TAG_CARD|TAG_HAS_CONTENT, TAG_DO|TAG_HAS_CONTENT,
			TAG_BR,
			END, END, END))
	// OUTPUT: <?xml version="1.0" encoding="utf-8"?>
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
	// OUTPUT: <?xml version="1.0" encoding="utf-8"?>
	// <XYZ xmlns="cp" xmlns:B="cp2"><B:CP2TAG><DO><BR/></DO></B:CP2TAG></XYZ>
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
	// OUTPUT: <?xml version="1.0" encoding="utf-8"?>
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
	// OUTPUT: <?xml version="1.0" encoding="utf-8"?>
	// <XYZ><CARD NAME="abc" STYLE="LIST"><DO TYPE="ACCEPT" URL="http://xyz.org/s"/> Enter name: <INPUT TYPE="TEXT" KEY="N"/></CARD></XYZ>
}
