package wbxml

import (
	"io"
	"testing"
)

func Test_ReadFromInsufficientBufferShouldReturnError(t *testing.T) {
	var st StringTable

	data := MakeDataBuffer(0x0A, 'H', 'e', 'l', 'l', 'o')

	err := st.ReadFromBuffer(data)
	if err == nil {
		t.Error("Error was nil but should be set")
	}

	if err != io.EOF {
		t.Errorf("Error should be io.EOF but was %s", err)
	}
}

func Test_GetString(t *testing.T) {
	var st StringTable

	data := MakeDataBuffer(
		0x0C,
		'H', 'e', 'l', 'l', 'o', 0x00,
		'W', 'o', 'r', 'l', 'd', 0x00)

	err := st.ReadFromBuffer(data)
	if err != nil {
		t.Errorf("Error should be nil but was %s", err)
	}

	value, err := st.getString(MakeDataBuffer(0x00))
	if err != nil {
		t.Errorf("Error should be nil but was %s", err)
	}

	if value != "Hello" {
		t.Errorf("Value should be 'Hello' but was '%s'", value)
	}

	value, err = st.getString(MakeDataBuffer(0x06))
	if err != nil {
		t.Errorf("Error should be nil but was %s", err)
	}

	if value != "World" {
		t.Errorf("Value should be 'World' but was '%s'", value)
	}
}
