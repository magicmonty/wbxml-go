package wbxml

import (
	"fmt"
	"testing"
)

func ExampleDecodeEntity() {
	decoder := NewDecoder(
		makeDataBuffer(
			ENTITY, 0x81, 0x20,
			ENTITY, 0x60),
		makeCodeBook())
	var result string
	result, _ = decoder.decodeEntity()
	fmt.Println(result)
	result, _ = decoder.decodeEntity()
	fmt.Println(result)
	// OUTPUT: &#160;
	// &#96;
}

func ExampleGetNameSpaceDeclarations() {
	decoder := NewDecoder(makeDataBuffer(0x00), makeCodeBook())

	decoder.usedNamespaces[0] = true
	decoder.usedNamespaces[1] = true
	decoder.usedNamespaces[255] = true
	fmt.Println(decoder.getNameSpaceDeclarations())
	// OUTPUT:  xmlns="cp" xmlns:B="cp2" xmlns:IV="cp255"
}

func TestGetNameSpaceDeclarationsShouldReturnEmptyStringIfOnlyCP0IsSelected(t *testing.T) {
	decoder := NewDecoder(makeDataBuffer(0x00), makeCodeBook())

	decoder.usedNamespaces[0] = true
	if decoder.getNameSpaceDeclarations() != "" {
		t.Errorf("NameSpace declaration should be emty but was \"%s\"", decoder.getNameSpaceDeclarations())
	}
}

func TestGetNameSpaceDeclarationsShouldReturnEmptyStringINoCPIsActive(t *testing.T) {
	decoder := NewDecoder(makeDataBuffer(0x00), makeCodeBook())

	if decoder.getNameSpaceDeclarations() != "" {
		t.Errorf("NameSpace declaration should be emty but was \"%s\"", decoder.getNameSpaceDeclarations())
	}
}
