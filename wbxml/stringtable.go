package wbxml

import (
	"bytes"
	"io"
)

type StringTable struct {
	length  uint32
	content []byte
}

func (st *StringTable) ReadFromBuffer(reader io.ByteReader) error {
	var err error

	st.length, err = readMultiByteUint32(reader)

	if err == nil && st.length > 0 {
		st.content = make([]byte, st.length)

		var index uint32
		for index = 0; index < st.length; index++ {
			st.content[index], err = reader.ReadByte()
			if err != nil {
				break
			}
		}
	}

	return err
}

func (st *StringTable) getString(reader io.ByteReader) (string, error) {
	var (
		result string = ""
		index  uint32
		b      bytes.Buffer
		err    error
	)

	index, err = readMultiByteUint32(reader)
	if err == nil {
		for i := index; i < uint32(len(st.content)); i++ {
			if st.content[i] != 0x00 {
				b.WriteByte(st.content[i])
			} else {
				break
			}
		}

		result, _ = b.ReadString(0x00)
	}

	return result, err
}
