package wbxml

import (
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
