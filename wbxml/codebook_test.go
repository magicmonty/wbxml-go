package wbxml

import (
	"testing"
)

func Test_NewCodeBook(t *testing.T) {
	var codeBook *CodeBook = NewCodeBook()

	if codeBook.HasCode(0) {
		t.Error("new code book should not have code page 0")
	}

	codeBook.AddCodePage(NewCodePage("Test", 0))

	if !codeBook.HasCode(0) {
		t.Error("code page 0 should exist")
	}

	if !codeBook.HasName("Test") {
		t.Error("name \"Test\" should exist")
	}

	codeBook.AddCodePage(NewCodePage("Test2", 1))

	if !codeBook.HasCode(1) {
		t.Error("code page 1 should exist after adding a new code page with code 1")
	}

	if !codeBook.HasName("Test2") {
		t.Error("name \"Test2\" should exist")
	}
}
