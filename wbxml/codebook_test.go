package wbxml

import (
	"testing"
)

func Test_NewCodeBook(t *testing.T) {
	var codeBook *CodeBook = NewCodeBook()

	if codeBook.IsReady() {
		t.Error("a new codebook should not be ready")
	}

	if codeBook.HasTagCode(0) {
		t.Error("new code book should not have tag code page 0")
	}

	if codeBook.HasAttributeCode(0) {
		t.Error("new code book should not have attribute code page 0")
	}

	codeBook.AddTagCodePage(NewCodePage("Test", 0))

	if !codeBook.HasTagCode(0) {
		t.Error("tag code page 0 should exist")
	}

	if !codeBook.IsReady() {
		t.Error("a codebook with code page 0 should be ready")
	}

	codeBook.AddTagCodePage(NewCodePage("Test2", 1))

	if !codeBook.HasTagCode(1) {
		t.Error("tag code page 1 should exist after adding a new tag code page with code 1")
	}

	codeBook.AddAttributeCodePage(NewAttributeCodePage(0))
	if !codeBook.HasAttributeCode(0) {
		t.Error("attribute code page 0 should exist")
	}

	codeBook.AddAttributeCodePage(NewAttributeCodePage(1))

	if !codeBook.HasAttributeCode(1) {
		t.Error("attribute code page 1 should exist after adding a new attribute code page with code 1")
	}

}
