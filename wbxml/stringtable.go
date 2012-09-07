package wbxml

import (
	"bytes"
)

type StringTable struct {
	length  byte
	content []byte
}

func (st *StringTable) ReadFromBuffer(b *bytes.Buffer) {
	st.length, _ = b.ReadByte()
	if st.length > 0 {
		st.content = make([]byte, st.length)

		var index byte
		for index = 0; index < st.length; index++ {
			st.content[index], _ = b.ReadByte()
		}
	}
}
