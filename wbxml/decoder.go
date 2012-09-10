package wbxml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

type Decoder struct {
	currentTagCodePage       CodePage
	currentAttributeCodePage AttributeCodePage
	usedNamespaces           map[byte]bool
	header                   Header
	reader                   io.ByteScanner
	codeBook                 *CodeBook
}

func NewDecoder(reader io.ByteScanner, codeBook *CodeBook) *Decoder {
	decoder := new(Decoder)
	decoder.codeBook = codeBook
	decoder.reader = reader
	decoder.usedNamespaces = make(map[byte]bool)

	return decoder
}

func (d *Decoder) decodeBody() (string, error) {
	var (
		documentType string = "<?xml version=\"1.0\"%s?>\n"
		result       string
		err          error
	)

	if d.codeBook.IsReady() {
		if d.header.charSetAsString == "" {
			documentType = fmt.Sprintf(documentType, "")
		} else {
			documentType = fmt.Sprintf(documentType, fmt.Sprintf(" encoding=\"%s\"", d.header.charSetAsString))
		}

		d.currentTagCodePage = d.codeBook.TagCodePages[0]
		if d.codeBook.HasAttributeCode(0) {
			d.currentAttributeCodePage = d.codeBook.AttributeCodePages[0]
		}

		result, err = d.decodeTag(true)
		if err == nil {
			result = documentType + result
		}
	} else {
		err = fmt.Errorf("CodeBook not ready")
	}

	return result, err
}

func (d *Decoder) decodeTag(addNamespaceDeclaration bool) (string, error) {
	var (
		nextByte byte
		err      error
	)

	nextByte, err = d.reader.ReadByte()
	if err == nil {
		d.reader.UnreadByte()

		if nextByte == SWITCH_PAGE {
			err = d.switchTagCodePage()
			if err == nil {
				return d.decodeTag(addNamespaceDeclaration)
			}
		} else {
			d.usedNamespaces[d.currentTagCodePage.Code] = true
			if nextByte&TAG_HAS_ATTRIBUTES != 0 {

				if nextByte&TAG_HAS_CONTENT != 0 {
					return d.decodeTagWithContentAndAttributes(addNamespaceDeclaration)
				} else {
					return d.decodeEmptyTagWithAttributes(addNamespaceDeclaration)
				}
			} else if nextByte&TAG_HAS_CONTENT != 0 {
				return d.decodeTagWithContent(addNamespaceDeclaration)
			} else {
				return d.decodeEmptyTag(addNamespaceDeclaration)
			}
		}
	}

	return "", err
}

func (d *Decoder) decodeTagWithContentAndAttributes(addNamespaceDeclaration bool) (string, error) {
	var (
		result     string = ""
		currentTag string
		content    string
		attributes string
		err        error
	)

	currentTag, err = d.decodeTagName()
	if err == nil && currentTag != "" {
		attributes, err = d.decodeAttributes()
		if err == nil {
			content, err = d.decodeContent()
			if err == nil {
				if addNamespaceDeclaration {
					result = fmt.Sprintf(
						"<%s%s%s\">%s</%s>",
						currentTag, d.getNameSpaceDeclarations(), attributes,
						content,
						currentTag)
				} else {
					result = fmt.Sprintf(
						"<%s%s\">%s</%s>",
						currentTag, attributes,
						content,
						currentTag)
				}
			}
		}
	}

	return result, err
}

func (d *Decoder) decodeEmptyTagWithAttributes(addNamespaceDeclaration bool) (string, error) {
	var (
		result     string = ""
		currentTag string
		attributes string
		err        error
	)
	currentTag, err = d.decodeTagName()
	if err == nil && currentTag != "" {
		attributes, err = d.decodeAttributes()
		if err == nil {
			if addNamespaceDeclaration {
				result = fmt.Sprintf(
					"<%s%s%s\"/>",
					currentTag, d.getNameSpaceDeclarations(), attributes)
			} else {
				result = fmt.Sprintf("<%s%s\"/>", currentTag, attributes)
			}
		}
	}

	return result, err
}

func (d *Decoder) decodeAttributes() (string, error) {
	var (
		result         string
		err            error
		nextByte       byte
		content        string
		firstAttribute bool = true
	)

	nextByte, err = d.reader.ReadByte()
	d.reader.UnreadByte()
	for nextByte != END {
		if err == nil {
			if nextByte == STR_T {
				content, err = d.decodeStringTableReference()
			} else if nextByte == STR_I {
				content, err = d.decodeInlineString()
			} else {
				content, err = d.decodeAttribute()
				if !firstAttribute && nextByte < ATTRIBUTE_VALUE_SPACE_START {
					content = "\"" + content
				}
				firstAttribute = false
			}

			if err == nil {
				result += content
			}
		} else {
			break
		}
		nextByte, err = d.reader.ReadByte()
		d.reader.UnreadByte()
	}

	if nextByte == END {
		d.reader.ReadByte()
	}

	return result, err
}

func (d *Decoder) decodeAttribute() (string, error) {
	var (
		nextByte byte
		err      error
		result   string
	)

	nextByte, err = d.reader.ReadByte()
	if err == nil {
		result = d.currentAttributeCodePage.GetString(nextByte)
	}

	return result, err
}

func (d *Decoder) decodeTagWithContent(addNamespaceDeclaration bool) (string, error) {
	var (
		result     string = ""
		nextByte   byte
		currentTag string
		content    string
		err        error
	)

	currentTag, err = d.decodeTagName()
	if err == nil && currentTag != "" {
		content, err = d.decodeContent()
		if err == nil {
			nextByte, err = d.reader.ReadByte()
			if err == nil && nextByte == END {
				if addNamespaceDeclaration {
					result = fmt.Sprintf(
						"<%s%s>%s</%s>",
						currentTag, d.getNameSpaceDeclarations(),
						content,
						currentTag)
				} else {
					result = fmt.Sprintf("<%s>%s</%s>", currentTag, content, currentTag)
				}
			}
		}
	}

	return result, err
}

func (d *Decoder) decodeContent() (string, error) {
	var (
		nextByte byte
		err      error
		result   string = ""
		content  string
	)

	nextByte, err = d.reader.ReadByte()
	d.reader.UnreadByte()
	for nextByte != END {
		if err == nil {

			if nextByte == STR_I {
				content, err = d.decodeInlineString()
			} else if nextByte == STR_T {
				content, err = d.decodeStringTableReference()
			} else if nextByte == SWITCH_PAGE {
				content = ""
				err = d.switchTagCodePage()
			} else if nextByte == ENTITY {
				content, err = d.decodeEntity()
			} else {
				content, err = d.decodeTag(false)
			}

			if err == nil {
				result += content
			}
		} else {
			break
		}
		nextByte, err = d.reader.ReadByte()
		d.reader.UnreadByte()
	}

	return result, err
}

func (d *Decoder) switchTagCodePage() error {
	var (
		nextByte byte
		err      error
	)

	nextByte, err = d.reader.ReadByte()
	if err == nil {
		if nextByte == SWITCH_PAGE {
			nextByte, err = d.reader.ReadByte()
			if err == nil {
				if d.codeBook.HasTagCode(nextByte) {
					d.usedNamespaces[nextByte] = true
					d.currentTagCodePage = d.codeBook.TagCodePages[nextByte]
				} else {
					err = fmt.Errorf("Codebook has no codepage %d", nextByte)
				}
			}
		} else {
			err = fmt.Errorf("Assumed SWITCH_PAGE token but was %d", nextByte)
		}
	}

	return err
}

func (d *Decoder) decodeInlineString() (string, error) {
	var (
		result   string = ""
		nextByte byte
		buffer   bytes.Buffer
		err      error
	)

	nextByte, err = d.reader.ReadByte()
	if err == nil && nextByte == STR_I {
		for true {
			nextByte, err = d.reader.ReadByte()
			if err != nil || nextByte == 0x00 {
				break
			}
			buffer.WriteByte(nextByte)
		}
		result, _ = buffer.ReadString(0x00)
	}

	return d.escapeString(result), err
}

func (d *Decoder) decodeEmptyTag(addNamespaceDeclaration bool) (string, error) {
	var (
		tagName string
		err     error
	)

	tagName, err = d.decodeTagName()

	if err == nil {
		if addNamespaceDeclaration {
			return "<" + tagName + d.getNameSpaceDeclarations() + "/>", nil
		} else {
			return "<" + tagName + "/>", nil
		}
	}

	return "", err
}

func (d *Decoder) decodeTagName() (string, error) {
	var (
		nextByte byte
		err      error
		tagName  string = ""
	)

	nextByte, err = d.reader.ReadByte()
	nextByte = nextByte &^ TAG_HAS_CONTENT &^ TAG_HAS_ATTRIBUTES
	if err == nil {
		if nextByte == LITERAL {
			tagName, err = d.header.stringTable.getString(d.reader)
			if err != nil {
				tagName = ""
			}
		} else if d.currentTagCodePage.HasTagCode(nextByte) {
			tagName = d.currentTagCodePage.Tags[nextByte]
		} else {
			err = fmt.Errorf("Unknown tag code: %d", nextByte)
		}

		if tagName != "" && d.currentTagCodePage.Code > 0 {
			tagName = d.currentTagCodePage.GetNameSpaceString() + ":" + tagName
		}
	}

	return tagName, err
}

func (d *Decoder) decodeEntity() (string, error) {
	var (
		result   string = ""
		entity   uint32
		nextByte byte
		err      error
	)

	nextByte, err = d.reader.ReadByte()
	if err == nil && nextByte == ENTITY {
		entity, err = readMultiByteUint32(d.reader)
		if err == nil {
			result = fmt.Sprintf("&#%d;", entity)
		}
	}

	return result, err
}

func (d *Decoder) escapeString(Value string) string {
	var b *bytes.Buffer = bytes.NewBuffer([]byte{})
	var writer io.Writer = b
	var result string

	xml.Escape(writer, []byte(Value))
	result, _ = b.ReadString(0x00)

	return result
}

func (d *Decoder) decodeStringTableReference() (string, error) {
	var (
		result   string = ""
		nextByte byte
		err      error
	)

	nextByte, err = d.reader.ReadByte()
	if err == nil {
		if nextByte == STR_T {
			return d.header.stringTable.getString(d.reader)
		}
	}

	return result, err
}

func (d *Decoder) getNameSpaceDeclarations() string {
	var (
		result                string
		i                     byte
		cp                    CodePage
		isOnlyCodePage0Active bool = true
	)

	for i = 0; i < 255; i++ {
		if d.usedNamespaces[i] {
			cp = d.codeBook.TagCodePages[i]
			result += cp.GetNameSpaceDeclaration()
			isOnlyCodePage0Active = i == 0
		}
	}

	if d.usedNamespaces[255] {
		cp = d.codeBook.TagCodePages[255]
		result += cp.GetNameSpaceDeclaration()
		isOnlyCodePage0Active = false
	}

	if isOnlyCodePage0Active {
		result = ""
	}

	return result
}
