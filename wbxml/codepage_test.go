package wbxml

import (
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

}
