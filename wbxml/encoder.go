package wbxml

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type Encoder struct {
	codeBook           *CodeBook
	xmlData            string
	writer             io.Writer
	currentTagCodePage CodePage
	stringTable        *StringTable

	rootElement    *element
	currentElement *element
}

func NewEncoder(codeBook *CodeBook, xmlData string, writer io.Writer) *Encoder {
	encoder := new(Encoder)
	encoder.codeBook = codeBook
	encoder.xmlData = xmlData
	encoder.writer = writer
	encoder.rootElement = nil
	encoder.currentElement = nil
	encoder.stringTable = NewStringTable()

	return encoder
}

func (e *Encoder) Encode() error {
	if e.codeBook.IsReady() {
		e.currentTagCodePage = e.codeBook.TagCodePages[0]

		r := strings.NewReader(e.xmlData)
		xmlDecoder := xml.NewDecoder(r)

		for t, err := xmlDecoder.Token(); err != io.EOF; t, err = xmlDecoder.Token() {
			if err != nil && err != io.EOF {
				return err
			}

			switch v := t.(type) {
			case xml.StartElement:
				err = e.encodeStartElement(interface{}(v).(xml.StartElement))
			case xml.EndElement:
				err = e.encodeEndElement(interface{}(v).(xml.EndElement))
			case xml.CharData:
				err = e.encodeCharData(interface{}(v).(xml.CharData))
			case xml.Comment:
				err = e.encodeComment(interface{}(v).(xml.Comment))
			case xml.ProcInst:
				err = e.encodeProcessingInstruction(interface{}(v).(xml.ProcInst))
			case xml.Directive:
				err = e.encodeDirective(interface{}(v).(xml.Directive))
			}

			if err != nil {
				return err
			}
		}

		if e.rootElement != nil {
			h := NewDefaultHeader()
			h.stringTable = e.stringTable

			err := h.Write(e.writer)
			if err != nil {
				return err
			}

			return e.rootElement.Encode(e.writer)
		}
	} else {
		return fmt.Errorf("CodeBook not ready")
	}

	return nil
}

func (e *Encoder) encodeStartElement(el xml.StartElement) error {
	tagName := el.Name.Local
	nameSpace := el.Name.Space

	if e.codeBook.HasNameSpace(nameSpace) || nameSpace == "" {
		if e.currentTagCodePage.Name != nameSpace {
			if nameSpace == "" {
				e.currentTagCodePage = e.codeBook.TagCodePages[0]
			} else {
				e.currentTagCodePage = e.codeBook.TagCodePagesByName[nameSpace]
			}
		}

		if e.currentTagCodePage.HasTag(tagName) {
			if e.rootElement == nil {
				e.rootElement = NewElement(
					nil,
					e.currentTagCodePage.TagCodes[tagName],
					e.currentTagCodePage.Code,
					e.encodeAttributes(el.Attr))
				e.currentElement = e.rootElement
			} else {
				var newElement = NewElement(
					e.currentElement,
					e.currentTagCodePage.TagCodes[tagName],
					e.currentTagCodePage.Code,
					e.encodeAttributes(el.Attr))
				e.currentElement.AddContent(newElement)
				e.currentElement = newElement
			}
		} else {
			var tagIndex uint32
			if e.stringTable == nil {
				e.stringTable = NewStringTable()
			}
			if e.stringTable.ContainsString(tagName) {
				tagIndex = e.stringTable.GetIndex(tagName)
			} else {
				tagIndex = e.stringTable.AddString(tagName)
			}

			if e.rootElement == nil {
				e.rootElement = NewElement(
					nil,
					LITERAL,
					e.currentTagCodePage.Code,
					e.encodeAttributes(el.Attr))
				e.rootElement.isLiteral = true
				e.rootElement.tagIndex = tagIndex
				e.currentElement = e.rootElement
			} else {
				var newElement = NewElement(
					e.currentElement,
					LITERAL,
					e.currentTagCodePage.Code,
					e.encodeAttributes(el.Attr))
				newElement.isLiteral = true
				newElement.tagIndex = tagIndex
				e.currentElement.AddContent(newElement)
				e.currentElement = newElement
			}
		}
	} else {
		return fmt.Errorf("NameSpace %s not found!", nameSpace)
	}

	return nil
}

func (e *Encoder) encodeAttributes(el []xml.Attr) []attribute {
	return nil
}

func (e *Encoder) encodeEndElement(el xml.EndElement) error {
	if e.currentElement != nil {
		e.currentElement = e.currentElement.parent
	}
	return nil
}

func (e *Encoder) encodeCharData(el xml.CharData) error {
	if e.currentElement != nil {
		if strings.Trim(fmt.Sprintf("%s", el), " \n\r\t") != "" {
			e.currentElement.AddContent(el.Copy())
		}
	}
	return nil
}

func (e *Encoder) encodeComment(el xml.Comment) error {
	return nil
}

func (e *Encoder) encodeProcessingInstruction(el xml.ProcInst) error {
	return nil
}

func (e *Encoder) encodeDirective(el xml.Directive) error {
	return nil
}
