package wbxml

import (
	"encoding/xml"
	"fmt"
	"io"
)

type contentElement interface{}
type attribute []byte

type element struct {
	tagCode    byte
	codePage   byte
	attributes []attribute
	content    []contentElement
	parent     *element
	tagIndex   uint32
	isLiteral  bool
}

func NewElement(parent *element, tagCode byte, codePage byte, attributes []attribute) *element {
	el := new(element)
	el.parent = parent
	el.tagCode = tagCode
	el.codePage = codePage
	el.attributes = attributes
	el.content = make([]contentElement, 0)
	el.tagIndex = 0
	el.isLiteral = false

	return el
}

func (e *element) HasContent() bool {
	return len(e.content) > 0
}

func (e *element) HasAttributes() bool {
	return len(e.attributes) > 0
}

func (e *element) IsRoot() bool {
	return e.parent == nil
}

func (e *element) AddContent(content contentElement) {
	e.content = append(e.content, content)
	el, ok := content.(*element)
	if ok {
		el.parent = e
	}
}

func (e *element) Encode(w io.Writer) error {
	return e.encodeTag(w, 0)
}

func (e *element) encodeTag(w io.Writer, currentCodePage byte) error {
	var (
		err error
		tag byte
	)

	if currentCodePage != e.codePage {
		_, err = w.Write([]byte{SWITCH_PAGE, e.codePage})
		if err != nil {
			return err
		}
	}

	tag = e.tagCode
	if e.HasContent() {
		tag |= TAG_HAS_CONTENT
	}

	if e.HasAttributes() {
		tag |= TAG_HAS_ATTRIBUTES
	}

	_, err = w.Write([]byte{tag})

	if e.isLiteral {
		err = writeMultiByteUint32(w, e.tagIndex)
		if err != nil {
			return err
		}
	}

	if e.HasAttributes() {
		err = e.encodeAttributes(w)
		if err != nil {
			return err
		}
	}

	if e.HasContent() {

		e.encodeContent(w)
		if err != nil {
			return err
		}

		_, err = w.Write([]byte{END})
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *element) encodeContent(w io.Writer) error {

	for i := 0; i < len(e.content); i++ {
		c := e.content[i]
		el, ok := c.(*element)
		if ok {
			el.encodeTag(w, e.codePage)
		} else {
			charData, ok := c.(xml.CharData)
			if ok {
				e.encodeCharData(w, charData)
			} else {
				return fmt.Errorf("Unknown element", charData)
			}
		}
	}

	return nil
}

func (e *element) encodeAttributes(w io.Writer) error {
	return nil
}

func (e *element) encodeCharData(w io.Writer, charData xml.CharData) error {
	first := true
	_, err := w.Write([]byte{STR_I})
	if err == nil {
		s := fmt.Sprintf("%s", charData)
		for _, c := range s {
			if c >= 0x20 && c <= 0x7E {
				_, err = w.Write([]byte{byte(c)})
				first = false
			} else {
				if !first {
					_, err = w.Write([]byte{0x00})
					if err != nil {
						return err
					}
				}

				first = false
				_, err = w.Write([]byte{ENTITY})
				if err == nil {
					err = writeMultiByteUint32(w, uint32(c))
					if err == nil {
						_, err = w.Write([]byte{STR_I})
					}
				}
			}
		}

		if err == nil {
			_, err = w.Write([]byte{0x00})
		}
	}
	return err
}
