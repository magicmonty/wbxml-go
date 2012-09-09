package wbxml

import (
	"fmt"
	"testing"
)

func CheckAttribute(t *testing.T, codePage AttributeCodePage, Name string, ValuePrefix string, Code byte) {
	if !codePage.HasCode(Code) {
		t.Errorf("The attribute for %d doesn't exist!", Code)
	}

	if codePage.Attributes[Code].Name != Name {
		t.Errorf("The name of the attribute for code %d should be %s but was %s", Code, Name, codePage.Attributes[Code].Name)
	}

	if codePage.Attributes[Code].ValuePrefix != ValuePrefix {
		t.Errorf("The value prefix of the attribute for code %d should be %s but was %s", Code, ValuePrefix, codePage.Attributes[Code].ValuePrefix)
	}
}

func CheckValue(t *testing.T, codePage AttributeCodePage, Value string, Code byte) {
	if !codePage.HasValueCode(Code) {
		t.Errorf("The attribute value for code %d doesn't exist!", Code)
	}

	if codePage.Values[Code] != Value {
		t.Errorf("The value for code %d should be %s but was %s", Code, Value, codePage.Values[Code])
	}
}

func Test_NewAttributeCodePage(t *testing.T) {
	var codePage AttributeCodePage = NewAttributeCodePage(0)
	codePage.AddAttribute("TYPE", "", 0x06)
	codePage.AddAttribute("TYPE", "TEXT", 0x07)
	codePage.AddAttribute("BOGUS", "", 0x87)

	CheckAttribute(t, codePage, "TYPE", "", 0x06)
	CheckAttribute(t, codePage, "TYPE", "TEXT", 0x07)

	if codePage.HasCode(0x87) {
		t.Error("The code 0x87 should not exist")
	}

	if codePage.HasCode(0x05) {
		t.Error("The code 0x05 should not exist")
	}

	if codePage.HasCode(0x00) {
		t.Error("The code 0x00 should not exist")
	}

	codePage.AddAttributeValue(".com", 0x0A)
	codePage.AddAttributeValue(".org", 0x85)
	codePage.AddAttributeValue("ACCEPT", 0x86)

	CheckValue(t, codePage, ".org", 0x85)
	CheckValue(t, codePage, "ACCEPT", 0x86)

	if codePage.HasValueCode(0x0A) {
		t.Error("The value code 0x0A should not exist")
	}

}

func ExampleToString() {
	var codePage AttributeCodePage = NewAttributeCodePage(0)
	codePage.AddAttribute("STYLE", "LIST", 0x05)
	codePage.AddAttribute("TYPE", "", 0x06)
	codePage.AddAttribute("TYPE", "TEXT", 0x07)
	codePage.AddAttribute("URL", "http://", 0x08)
	codePage.AddAttribute("NAME", "", 0x09)
	codePage.AddAttribute("KEY", "", 0x0A)

	codePage.AddAttributeValue(".org", 0x85)
	codePage.AddAttributeValue("ACCEPT", 0x86)

	fmt.Println(codePage.GetString(0x05))
	fmt.Println(codePage.GetString(0x06))
	fmt.Println(codePage.GetString(0x07))
	fmt.Println(codePage.GetString(0x08))
	fmt.Println(codePage.GetString(0x09))
	fmt.Println(codePage.GetString(0x0A))
	fmt.Println(codePage.GetString(0x0C)) // code doesn't existent 
	fmt.Println(codePage.GetString(0x85))
	fmt.Println(codePage.GetString(0x86))
	// OUTPUT: STYLE="LIST
	// TYPE=
	// TYPE="TEXT
	// URL="http://
	// NAME=
	// KEY=
	//
	// .org
	// ACCEPT
}
