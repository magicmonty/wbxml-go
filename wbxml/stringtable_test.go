package wbxml

import (
	"io"
	"testing"
)

func Test_ReadFromInsufficientBufferShouldReturnError(t *testing.T) {
	var st StringTable

	data := makeDataBuffer(0x0A, 'H', 'e', 'l', 'l', 'o')

	err := st.Read(data)
	if err == nil {
		t.Error("Error was nil but should be set")
	}

	if err != io.EOF {
		t.Errorf("Error should be io.EOF but was %s", err)
	}
}

func Test_GetString(t *testing.T) {
	var st StringTable

	data := makeDataBuffer(
		0x0C,
		'H', 'e', 'l', 'l', 'o', 0x00,
		'W', 'o', 'r', 'l', 'd', 0x00)

	err := st.Read(data)
	if err != nil {
		t.Errorf("Error should be nil but was %s", err)
	}

	value, err := st.getString(makeDataBuffer(0x00))
	if err != nil {
		t.Errorf("Error should be nil but was %s", err)
	}

	if value != "Hello" {
		t.Errorf("Value should be 'Hello' but was '%s'", value)
	}

	value, err = st.getString(makeDataBuffer(0x06))
	if err != nil {
		t.Errorf("Error should be nil but was %s", err)
	}

	if value != "World" {
		t.Errorf("Value should be 'World' but was '%s'", value)
	}
}

func Test_ContainsString(t *testing.T) {
	st := NewStringTable()

	data := makeDataBuffer(
		0x0C,
		'H', 'e', 'l', 'l', 'o', 0x00,
		'W', 'o', 'r', 'l', 'd', 0x00)
	st.Read(data)

	if !st.ContainsString("Hello") {
		t.Error("string table should contain string 'Hello'")
	}

	if !st.ContainsString("World") {
		t.Error("string table should contain string 'World'")
	}
}

func Test_GetIndex(t *testing.T) {
	st := NewStringTable()

	data := makeDataBuffer(
		0x0C,
		'H', 'e', 'l', 'l', 'o', 0x00,
		'W', 'o', 'r', 'l', 'd', 0x00)
	st.Read(data)

	if st.GetIndex("Hello") != 0 {
		t.Errorf("index of string 'Hello' should be 0 but was %d", st.GetIndex("Hello"))
	}

	if st.GetIndex("World") != 6 {
		t.Errorf("index of string 'World' should be 0 but was %d", st.GetIndex("World"))
	}
}

func Test_AddString(t *testing.T) {
	st := NewStringTable()

	if st.AddString("Hello") != 0 {
		t.Errorf("AddString('Hello') should return 0")
	}

	if st.AddString("World") != 6 {
		t.Errorf("AddString('World') should return 6")
	}

	if st.length != 12 {
		t.Errorf("Length of string buffer should be 12 but is %d", st.length)
	}

	if st.ContainsString("Hello") {
		if st.GetIndex("Hello") != 0 {
			t.Errorf("index of string 'Hello' should be 0 but was %d", st.GetIndex("Hello"))
		}
	} else {
		t.Error("string table should contain string 'Hello'")
	}

	if st.ContainsString("World") {
		if st.GetIndex("World") != 6 {
			t.Errorf("index of string 'World' should be 6 but was %d", st.GetIndex("World"))
		}
	} else {
		t.Error("string table should contain string 'World'")
	}
}
