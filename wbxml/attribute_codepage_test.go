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
	// OUTPUT:  STYLE="LIST
	//  TYPE="
	//  TYPE="TEXT
	//  URL="http://
	//  NAME="
	//  KEY="
	//
	// .org
	// ACCEPT
}

func Test_HasAttribute(t *testing.T) {
	var codePage AttributeCodePage = NewAttributeCodePage(0)
	codePage.AddAttribute("STYLE", "LIST", 0x05)
	codePage.AddAttribute("TYPE", "", 0x06)
	codePage.AddAttribute("TYPE", "TEXT", 0x07)
	codePage.AddAttribute("URL", "http://", 0x08)
	codePage.AddAttribute("NAME", "", 0x09)
	codePage.AddAttribute("KEY", "", 0x0A)

	if codePage.HasAttribute("STYLE", "") {
		t.Error("codePage should not have attribute STYLE")
	}

	if !codePage.HasAttribute("STYLE", "LIST") {
		t.Error("codePage should have attribute STYLE=LIST")
	}

	if !codePage.HasAttribute("TYPE", "") {
		t.Error("codePage should have attribute TYPE")
	}

	if !codePage.HasAttribute("TYPE", "TEXT") {
		t.Error("codePage should have attribute TYPE=TEXT")
	}

	if !codePage.HasAttribute("TYPE", "PASSWORD") {
		t.Error("codePage should have attribute TYPE=PASSWORD")
	}

	if !codePage.HasAttribute("NAME", "") {
		t.Error("codePage should have attribute NAME")
	}

	if !codePage.HasAttribute("NAME", "Karl") {
		t.Error("codePage should have attribute NAME")
	}

	if !codePage.HasAttribute("KEY", "") {
		t.Error("codePage should have attribute KEY")
	}

	if !codePage.HasAttribute("URL", "http://www.google.de") {
		t.Error("codePage should have attribute URL=www.google.de")
	}
}

func checkAttributeToken(expected byte, actual byte, t *testing.T) {
	if actual != expected {
		t.Errorf("attribute token should be 0x%0.2X but was 0x%0.2X", expected, actual)
	}
}

func checkValueToken(t *testing.T, actual []string, expected ...string) {
	if actual == nil {
		t.Errorf("value tokens should not be nil")
	}

	if len(actual) != len(expected) {
		t.Errorf("Length of vt should be %d but was %d (expected: %s, actual: %s)", len(expected), len(actual), expected, actual)
	} else {
		for i := 0; i < len(expected); i++ {
			if actual[i] != expected[i] {
				t.Errorf("value token %d should be %s but was '%s'", i, expected[i], actual[i])
			}
		}
	}
}

func Test_Tokenize(t *testing.T) {
	var codePage AttributeCodePage = NewAttributeCodePage(0)
	codePage.AddAttribute("STYLE", "LIST", 0x05)
	codePage.AddAttribute("TYPE", "", 0x06)
	codePage.AddAttribute("TYPE", "TEXT", 0x07)
	codePage.AddAttribute("URL", "http://", 0x08)
	codePage.AddAttribute("NAME", "", 0x09)
	codePage.AddAttribute("KEY", "", 0x0A)

	codePage.AddAttributeValue(".org", 0x85)
	codePage.AddAttributeValue("ACCEPT", 0x86)
	codePage.AddAttributeValue(".com", 0x8A)
	codePage.AddAttributeValue("google", 0x8B)

	at, vt, err := codePage.Tokenize("STYLE", "")
	if err == nil || err.Error() != "codepage has no matching attribute entry (name='STYLE', value='')" {
		t.Error("unsupported Attribute should trigger error")
	}

	checkAttributeToken(0x00, at, t)
	if vt != nil {
		t.Errorf("Invalid value tokens should be nil but was ", vt)
	}

	at, vt, err = codePage.Tokenize("URL", "http://www.google.com")
	checkNoError(t, err)
	checkAttributeToken(0x08, at, t)
	checkValueToken(t, vt, "www.", "google", ".com")

	at, vt, err = codePage.Tokenize("URL", "http://www.google.de")
	checkNoError(t, err)
	checkAttributeToken(0x08, at, t)
	checkValueToken(t, vt, "www.", "google", ".de")

	at, vt, err = codePage.Tokenize("URL", "http://www.suse.de")
	checkNoError(t, err)
	checkAttributeToken(0x08, at, t)
	checkValueToken(t, vt, "www.suse.de")

	at, vt, err = codePage.Tokenize("TYPE", "BUTTON")
	checkNoError(t, err)
	checkAttributeToken(0x06, at, t)
	checkValueToken(t, vt, "BUTTON")

	at, vt, err = codePage.Tokenize("TYPE", "TEXT")
	checkNoError(t, err)
	checkAttributeToken(0x07, at, t)
	checkValueToken(t, vt)

	at, vt, err = codePage.Tokenize("TYPE", "TEXT1")
	checkNoError(t, err)
	checkValueToken(t, vt, "1")
}
