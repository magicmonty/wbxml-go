package wbxml

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type Encoder struct {
	codeBook                 *CodeBook
	xmlData                  string
	writer                   io.Writer
	currentTagCodePage       CodePage
	currentAttributeCodePage AttributeCodePage
	stringTable              *StringTable

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
		if e.codeBook.AttributeCodePages != nil {
			e.currentAttributeCodePage = e.codeBook.AttributeCodePages[0]
		}

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
		e.updateCurrentCodePage(nameSpace)
		newElement := e.makeNewElement(tagName, el.Attr)
		e.updateCurrentElement(newElement)
	} else {
		return fmt.Errorf("NameSpace %s not found!", nameSpace)
	}

	return nil
}

func (e *Encoder) updateCurrentCodePage(nameSpace string) {
	if e.currentTagCodePage.Name != nameSpace {
		if nameSpace == "" {
			e.currentTagCodePage = e.codeBook.TagCodePages[0]
		} else {
			e.currentTagCodePage = e.codeBook.TagCodePagesByName[nameSpace]
		}
	}
}

func (e *Encoder) makeNewElement(tagName string, attributes []xml.Attr) *element {
	var newElement *element = nil
	if e.currentTagCodePage.HasTag(tagName) {
		newElement = e.makeNewElementFromTagCode(
			e.currentTagCodePage.TagCodes[tagName],
			attributes)
	} else {
		newElement = e.makeNewLiteralElement(
			tagName,
			attributes)
	}

	return newElement
}

func (e *Encoder) makeNewLiteralElement(tagName string, attributes []xml.Attr) *element {
	var tagIndex uint32
	if e.stringTable == nil {
		e.stringTable = NewStringTable()
	}

	if e.stringTable.ContainsString(tagName) {
		tagIndex = e.stringTable.GetIndex(tagName)
	} else {
		tagIndex = e.stringTable.AddString(tagName)
	}

	newElement := e.makeNewElementFromTagCode(
		LITERAL, attributes)
	newElement.isLiteral = true
	newElement.tagIndex = tagIndex

	return newElement
}

func (e *Encoder) makeNewElementFromTagCode(code byte, attributes []xml.Attr) *element {
	attrs, err := e.encodeAttributes(attributes)
	if err == nil {
		return NewElement(code, e.currentTagCodePage.Code, attrs)
	}

	return nil
}

func (e *Encoder) updateCurrentElement(newElement *element) {
	if e.rootElement == nil {
		e.rootElement = newElement
		e.currentElement = e.rootElement
	} else {
		e.currentElement.AddContent(newElement)
		e.currentElement = newElement
	}
}

func (e *Encoder) encodeAttributes(el []xml.Attr) ([]attribute, error) {
	if e.codeBook.AttributeCodePages != nil {
		result := make([]attribute, 0)
		for _, a := range el {
			if !strings.Contains(a.Name.Local, "xmlns") {
				if e.currentAttributeCodePage.HasAttribute(a.Name.Local, a.Value) {
					at, vt, err := e.currentAttributeCodePage.Tokenize(a.Name.Local, a.Value)
					if err == nil {
						attr := make(attribute, 0)
						attr = append(attr, at)

						for _, v := range vt {
							if e.currentAttributeCodePage.HasValue(v) {
								attr = append(attr, e.currentAttributeCodePage.ValueCodes[v])
							} else {
								attr = append(attr, STR_I)
								for _, c := range v {
									attr = append(attr, byte(c))
								}
								attr = append(attr, 0x00)
							}
						}
						result = append(result, attr)
					} else {
						return nil, err
					}
				}
			}
		}
		return result, nil
	}

	return nil, fmt.Errorf("CodeBook has no attribute codepages!")
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
