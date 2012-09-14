package wbxml

import (
	"bytes"
	"fmt"
)

func encodeExample(xml string) {
	w := bytes.NewBuffer(make([]byte, 0))
	e := NewEncoder(
		makeCodeBook(),
		xml,
		w)
	err := e.Encode()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		printByteStream(w)
	}
}

// TODO
func ExampleEncodeWbxml() {
	encodeExample(
		`<?xml version="1.0" encoding="utf-8"?>
			<XYZ xmlns="cp" xmlns:B="cp2">
				<CARD>
					<B:CP2TAG>
						<DO>Hello&#160;World</DO>
					</B:CP2TAG>
				</CARD>
			</XYZ>`)
	// OUTPUT: 03 01 6A 00 47 46 00 01 45 00 00 48 03 48 65 6C 6C 6F 00 02 81 20 03 57 6F 72 6C 64 00 01 01 01 01 
}

func ExampleEncodeEmptyTag() {
	encodeExample(
		`<?xml version="1.0" encoding="utf-8"?>
		<XYZ/>`)
	// OUTPUT: 03 01 6A 00 07
}

func ExampleEncodeEmptyTagWithDifferentNameSpace() {
	encodeExample(
		`<?xml version="1.0" encoding="utf-8"?>
		<B:CP2TAG xmlns:B="cp2"/>`)
	// OUTPUT: 03 01 6A 00 00 01 05
}

func ExampleEncodeEmptyLiteralTag() {
	encodeExample(
		`<?xml version="1.0" encoding="utf-8"?>
		<ABC/>`)
	// OUTPUT: 03 01 6A 04 41 42 43 00 04 00
}

func ExampleEncodeTagWithEmptyTagAsContent() {
	encodeExample(
		`<?xml version="1.0" encoding="utf-8"?>
	    <XYZ><CARD/></XYZ>`)
	// OUTPUT: 03 01 6A 00 47 06 01
}

func ExampleEncodeTagWithEmptyTagFromDifferentCodePageAsContent() {
	encodeExample(
		`<?xml version="1.0" encoding="utf-8"?>
	     <XYZ xmlns="cp" xmlns:B="cp2">
	     	<B:CP2TAG/>
	     </XYZ>`)
	// OUTPUT: 03 01 6A 00 47 00 01 05 01
}

func ExampleEncodeTagWithTextAsContent() {
	encodeExample(
		`<?xml version="1.0" encoding="utf-8"?>
	     <XYZ>X &amp; Y</XYZ>`)
	// OUTPUT: 03 01 6A 00 47 03 58 20 26 20 59 00 01
}

func ExampleEncodeTagFromDifferentCodePageWithTextAsContent() {
	encodeExample(
		`<?xml version="1.0" encoding="utf-8"?>
		<B:CP2TAG xmlns:B="cp2">X &amp; Y</B:CP2TAG>`)
	// OUTPUT: 03 01 6A 00 00 01 45 03 58 20 26 20 59 00 01
}

func ExampleEncodeMultipleNestedTags() {
	encodeExample(
		`<?xml version="1.0" encoding="utf-8"?>
		<XYZ>
			<CARD>
				<DO>
					<BR/>
				</DO>
			</CARD>
		</XYZ>`)
	// OUTPUT: 03 01 6A 00 47 46 48 05 01 01 01
}

func ExampleEncodeMultipleNestedTagsWithDifferentCodePages() {
	encodeExample(
		`<?xml version="1.0" encoding="utf-8"?>
		<XYZ xmlns="cp" xmlns:B="cp2">
			<B:CP2TAG>
				<DO>
					<BR/>
				</DO>
			</B:CP2TAG>
		</XYZ>`)
	// OUTPUT: 03 01 6A 00 47 00 01 45 00 00 48 05 01 01 01
}

// Example from http://www.w3.org/TR/wbxml/#_Toc443384926
// TODO
func ExampleSimpleWBXMLEncode() {
	encodeExample(
		`<?xml version="1.0" encoding="utf-8"?>
		<XYZ><CARD> X &amp; Y<BR/> X&#160;=&#160;1 </CARD></XYZ>`)
	// OUTPUT: 03 01 6A 00 47 46 03 20 58 20 26 20 59 00 05 03 20 58 00 02 81 20 03 3D 00 02 81 20 03 31 20 00 01 01
}

// Example from http://www.w3.org/TR/wbxml/#_Toc443384927
// TODO
func _ExampleExtendedWBXMLEncode() {
	encodeExample(
		`<?xml version="1.0" encoding="utf-8"?>
	    <XYZ>
	    	<CARD NAME="abc" STYLE="LIST">
	    		<DO TYPE="ACCEPT" URL="http://xyz.org/s"/> Enter name: <INPUT TYPE="TEXT" KEY="N"/>
	    	</CARD>
	    </XYZ>`)
	// OUTPUT: 03 01 6A 12 61 62 63 00 20 45 6E 74 65 72 20 6E 61 6D 65 3A 20 00 47 C6 09 83 00 05 01 88 06 86 08 03 78 79 7A 00 85 03 2F 73 00 01 83 04 89 07 0A 03 4E 00 01 01 01
}
