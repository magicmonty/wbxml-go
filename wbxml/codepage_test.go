package wbxml

import (
	"fmt"
	"testing"
)

func CheckTag(t *testing.T, codePage CodePage, Tag string, Code byte) {
	if !codePage.HasTag(Tag) {
		t.Errorf("The tag %s doesn't exist!", Tag)
	}

	if codePage.Tags[Code] != Tag {
		t.Errorf("The tag for %d should be %s", Code, Tag)
	}

	if codePage.TagCodes[Tag] != Code {
		t.Errorf("The code for %s should be %d", Tag, Code)
	}
}

func Test_NewCodePage(t *testing.T) {
	var codePage CodePage = NewCodePage("Test", 0)
	codePage.AddTag("Sync", 0x05)
	codePage.AddTag("Responses", 0x06)

	CheckTag(t, codePage, "Sync", 0x05)
	CheckTag(t, codePage, "Responses", 0x06)

	if codePage.HasTag("Blubb") {
		t.Error("The tag \"Blubb\" should not exist")
	}

	if codePage.HasTagCode(0x08) {
		t.Error("The tag code 0x08 should not exist")
	}

	if codePage.HasTagCode(0x00) {
		t.Error("The tag code 0x00 should not exist")
	}
}

func Example_GetNameSpaceString() {
	var (
		codePage0   CodePage = NewCodePage("cp0", 0x00)
		codePage1   CodePage = NewCodePage("cp1", 0x01)
		codePage25  CodePage = NewCodePage("cp25", 0x19)
		codePage26  CodePage = NewCodePage("cp26", 0x1A)
		codePage51  CodePage = NewCodePage("cp51", 0x33)
		codePage52  CodePage = NewCodePage("cp52", 0x34)
		codePage255 CodePage = NewCodePage("cp255", 0xFF)
	)

	fmt.Println(codePage0.GetNameSpaceString())
	fmt.Println(codePage1.GetNameSpaceString())
	fmt.Println(codePage25.GetNameSpaceString())
	fmt.Println(codePage26.GetNameSpaceString())
	fmt.Println(codePage51.GetNameSpaceString())
	fmt.Println(codePage52.GetNameSpaceString())
	fmt.Println(codePage255.GetNameSpaceString())
	// OUTPUT: 
	// B
	// Z
	// AA
	// AZ
	// BA
	// IV
}

func Example_GetNameSpaceDeclaration() {
	var (
		codePage0   CodePage = NewCodePage("cp0", 0x00)
		codePage1   CodePage = NewCodePage("cp1", 0x01)
		codePage25  CodePage = NewCodePage("cp25", 0x19)
		codePage26  CodePage = NewCodePage("cp26", 0x1A)
		codePage51  CodePage = NewCodePage("cp51", 0x33)
		codePage52  CodePage = NewCodePage("cp52", 0x34)
		codePage255 CodePage = NewCodePage("cp255", 0xFF)
	)

	fmt.Println(codePage0.GetNameSpaceDeclaration())
	fmt.Println(codePage1.GetNameSpaceDeclaration())
	fmt.Println(codePage25.GetNameSpaceDeclaration())
	fmt.Println(codePage26.GetNameSpaceDeclaration())
	fmt.Println(codePage51.GetNameSpaceDeclaration())
	fmt.Println(codePage52.GetNameSpaceDeclaration())
	fmt.Println(codePage255.GetNameSpaceDeclaration())
	// OUTPUT:  xmlns="cp0"
	//  xmlns:B="cp1"
	//  xmlns:Z="cp25"
	//  xmlns:AA="cp26"
	//  xmlns:AZ="cp51"
	//  xmlns:BA="cp52"
	//  xmlns:IV="cp255"
}
