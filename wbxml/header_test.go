package wbxml

import (
	"testing"
)

func Test_HeaderCharset(t *testing.T) {
	var (
		h   *Header = NewDefaultHeader()
		err error
	)

	err = h.Read(makeDataBuffer(WBXML_1_3, UNKNOWN_PI, CHARSET_UTF8, 0x00))

	if err != nil {
		t.Errorf("Error should be nil but was %s", err.Error())
	}

	if h.charSetAsString != "utf-8" {
		t.Errorf("Charset should be \"utf-8\" but was \"%s\"", h.charSetAsString)
	}

	err = h.Read(makeDataBuffer(WBXML_1_3, UNKNOWN_PI, 111, 0x00))

	if err != nil {
		t.Errorf("Error should be nil but was %s", err.Error())
	}

	if h.charSetAsString != "iso-8859-15" {
		t.Errorf("Charset should be \"iso-8859-15\" but was \"%s\"", h.charSetAsString)
	}

	err = h.Read(makeDataBuffer(WBXML_1_3, UNKNOWN_PI, 0x00, 0x00))

	if h.charSetAsString != "" {
		t.Errorf("Charset should be empty but was \"%s\"", h.charSetAsString)
	}
}
