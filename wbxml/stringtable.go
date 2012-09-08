package wbxml

import (
	"bytes"
	"io"
)

type StringTable struct {
	length  uint32
	content []byte
}

func (st *StringTable) ReadFromBuffer(reader io.ByteReader) {
	st.length = readMultiByteUint32(reader)
	if st.length > 0 {
		st.content = make([]byte, st.length)

		var index uint32
		for index = 0; index < st.length; index++ {
			st.content[index], _ = reader.ReadByte()
		}
	}
}

func (st *StringTable) getString(reader io.ByteReader) string {
	var (
		result string = ""
		index  uint32
		b      bytes.Buffer
	)

	index = readMultiByteUint32(reader)
	for i := index; i < uint32(len(st.content)); i++ {
		if st.content[i] != 0x00 {
			b.WriteByte(st.content[i])
		} else {
			break
		}
	}

	result, _ = b.ReadString(0x00)
	return result
}
