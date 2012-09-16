package wbxml

import (
	"testing"
)

func Test_NewElementShouldHaveNoContent(t *testing.T) {
	el := NewElement(5, 0, nil)

	if el.HasContent() {
		t.Error("Element should have no content")
	}
}

func Test_ElementWithoutAttributesShouldHaveNoContent(t *testing.T) {
	el := NewElement(5, 0, nil)
	if el.HasAttributes() {
		t.Error("Element should have no attributes")
	}
}

func Test_AddContent(t *testing.T) {
	el := NewElement(5, 0, nil)
	el2 := NewElement(6, 0, nil)

	el.AddContent(el2)

	if !el.HasContent() {
		t.Error("Element should have content")
	}

	if el2.parent != el {
		t.Error("Inner element should have outer element as parent")
	}
}
