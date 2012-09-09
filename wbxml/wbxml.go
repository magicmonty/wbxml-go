package wbxml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

const (
	TAG_STATE       byte = 1
	ATTRIBUTE_STATE byte = 2
)

type Decoder struct {
	currentTagCodePage       CodePage
	currentAttributeCodePage AttributeCodePage
	currentState             byte
	header                   Header
	reader                   io.ByteScanner
	codeBook                 *CodeBook
}

func NewDecoder(reader io.ByteScanner, codeBook *CodeBook) (*Decoder, error) {
	var err error

	decoder := new(Decoder)
	decoder.codeBook = codeBook
	decoder.reader = reader
	decoder.currentTagCodePage = codeBook.TagCodePages[0]
	if codeBook.HasAttributeCode(0) {
		decoder.currentAttributeCodePage = codeBook.AttributeCodePages[0]
	}
	decoder.currentState = TAG_STATE
	err = decoder.header.ReadFromBuffer(reader)
	if err != nil {
		decoder = nil
	}
	return decoder, err
}

func Decode(reader io.ByteScanner, codeBook *CodeBook) (string, error) {
	if codeBook.IsReady() {
		decoder, err := NewDecoder(reader, codeBook)
		if err == nil {
			return decoder.decodeBody()
		} else {
			return "", err
		}

	}

	return "", nil
}

func (d *Decoder) decodeBody() (string, error) {
	var (
		documentType string = "<?xml version=\"1.0\"?>\n"
		result       string
		err          error
	)

	result, err = d.decodeTag()
	if err == nil {
		result = documentType + result
	}

	return result, err
}

func (d *Decoder) decodeTag() (string, error) {
	var (
		nextByte byte
		err      error
	)

	nextByte, err = d.reader.ReadByte()
	if err == nil {
		d.reader.UnreadByte()

		if nextByte&TAG_HAS_ATTRIBUTES != 0 {

			if nextByte&TAG_HAS_CONTENT != 0 {
				return d.decodeTagWithContentAndAttributes()
			} else {
				return d.decodeEmptyTagWithAttributes()
			}
		} else if nextByte&TAG_HAS_CONTENT != 0 {
			return d.decodeTagWithContent()
		} else {
			return d.decodeEmptyTag()
		}
	}

	return "", err
}

func (d *Decoder) decodeTagWithContentAndAttributes() (string, error) {
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
				result = fmt.Sprintf("<%s%s\">%s</%s>", currentTag, attributes, content, currentTag)
			}
		}
	}

	return result, err
}

func (d *Decoder) decodeEmptyTagWithAttributes() (string, error) {
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
			result = fmt.Sprintf("<%s%s\"/>", currentTag, attributes)
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

func (d *Decoder) decodeTagWithContent() (string, error) {
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
				result = fmt.Sprintf("<%s>%s</%s>", currentTag, content, currentTag)
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
			} else if nextByte == ENTITY {
				content, err = d.decodeEntity()
			} else {
				content, err = d.decodeTag()
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

func (d *Decoder) decodeEmptyTag() (string, error) {
	var (
		tagName string
		err     error
	)

	tagName, err = d.decodeTagName()

	if err == nil {
		return "<" + tagName + "/>", nil
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
