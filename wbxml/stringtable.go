package wbxml

import (
	"io"
)

type StringTable struct {
	length  byte
	content []byte
}

func (st *StringTable) ReadFromBuffer(reader io.ByteReader) {
	st.length, _ = reader.ReadByte()
	if st.length > 0 {
		st.content = make([]byte, st.length)

		var index byte
		for index = 0; index < st.length; index++ {
			st.content[index], _ = reader.ReadByte()
		}
	}
}
