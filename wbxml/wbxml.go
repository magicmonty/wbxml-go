package wbxml

import (
	"bytes"
	"fmt"
	"io"
)

const (
	TAG_STATE       byte = 1
	ATTRIBUTE_STATE byte = 2
)

type Decoder struct {
	currentTagCodePage       CodePage
	currentAttributeCodePage CodePage
	currentState             byte
	header                   Header
	reader                   io.ByteScanner
	codeBook                 *CodeBook
}

func NewDecoder(reader io.ByteScanner, codeBook *CodeBook) *Decoder {
	decoder := new(Decoder)
	decoder.codeBook = codeBook
	decoder.reader = reader
	decoder.currentTagCodePage = codeBook.CodePages[0]
	decoder.currentState = TAG_STATE
	decoder.header.ReadFromBuffer(reader)
	return decoder
}

func (d *Decoder) decodeEntity() string {
	var (
		result   string = ""
		entity   uint32
		nextByte byte
	)

	nextByte, _ = d.reader.ReadByte()
	if nextByte == ENTITY {
		entity = readMultiByteUint32(d.reader)
		result = fmt.Sprintf("&#%d;", entity)
	}
	return result
}

func (d *Decoder) decodeInlineString() string {
	var (
		result   string = ""
		nextByte byte
		buffer   bytes.Buffer
	)

	nextByte, _ = d.reader.ReadByte()
	if nextByte == STR_I {
		for true {
			nextByte, _ = d.reader.ReadByte()
			if nextByte == 0x00 {
				break
			}
			buffer.WriteByte(nextByte)
		}
		result, _ = buffer.ReadString(0x00)
	}
	return result
}

func (d *Decoder) decodeStringTableReference() string {
	var (
		result   string = ""
		nextByte byte
	)

	nextByte, _ = d.reader.ReadByte()
	if nextByte == STR_T {
	}

	return result
}

func (d *Decoder) decodeTagWithContentAndAttributes() string {
	return ""
}

func (d *Decoder) decodeTagWithContent() string {
	var (
		result     string = ""
		nextByte   byte
		tagCode    byte
		currentTag string
	)

	nextByte, _ = d.reader.ReadByte()
	tagCode = nextByte &^ TAG_HAS_CONTENT
	if d.currentTagCodePage.HasTagCode(tagCode) {
		currentTag = d.currentTagCodePage.Tags[tagCode]
		result = "<" + currentTag + ">"
		result += d.decodeElement()
		nextByte, _ = d.reader.ReadByte()
		if nextByte == END {
			result += "</" + currentTag + ">"
		}
	}

	return result
}

func (d *Decoder) decodeEmptyTagWithAttributes() string {
	return ""
}

func (d *Decoder) decodeEmptyTag() string {
	var (
		nextByte byte
	)

	nextByte, _ = d.reader.ReadByte()

	if nextByte == LITERAL {
		return "<" + d.header.stringTable.getString(d.reader) + "/>"
	} else if d.currentTagCodePage.HasTagCode(nextByte) {
		return "<" + d.currentTagCodePage.Tags[nextByte] + "/>"
	}

	return ""
}

func (d *Decoder) decodeElement() string {
	var (
		nextByte byte
	)

	nextByte, _ = d.reader.ReadByte()
	d.reader.UnreadByte()

	if nextByte == STR_I {
		return d.decodeInlineString()
	} else if nextByte&TAG_HAS_CONTENT != 0 {

		if nextByte&TAG_HAS_ATTRIBUTES != 0 {
			return d.decodeTagWithContentAndAttributes()
		} else {
			return d.decodeTagWithContent()
		}
	} else if nextByte&TAG_HAS_ATTRIBUTES != 0 {
		return d.decodeEmptyTagWithAttributes()
	} else {
		return d.decodeEmptyTag()
	}

	return ""
}

func (d *Decoder) decodeBody() string {
	var (
		documentType string = "<?xml version=\"1.0\"?>\n"
	)

	return documentType + d.decodeElement()
}

func Decode(reader io.ByteScanner, codeBook *CodeBook) string {
	return NewDecoder(reader, codeBook).decodeBody()
}
