package wbxml

import (
	"bytes"
	"fmt"
	"io"
)

type StringTable struct {
	length  uint32
	content []byte
	index   map[string]uint32
}

func NewStringTable() *StringTable {
	s := new(StringTable)

	s.length = 0
	s.content = make([]byte, 0)
	s.index = make(map[string]uint32, 0)

	return s
}

func (st *StringTable) Read(reader io.ByteReader) error {
	var (
		err           error
		currentString string = ""
		currentIndex  uint32 = 0
		b             byte
	)

	st.index = make(map[string]uint32)
	st.length, err = readMultiByteUint32(reader)

	if err == nil && st.length > 0 {
		st.content = make([]byte, st.length)

		var index uint32
		for index = 0; index < st.length; index++ {
			b, err = reader.ReadByte()
			if err != nil {
				break
			}

			st.content[index] = b
			if b > 0 {
				currentString += fmt.Sprintf("%c", b)
			} else {
				st.index[currentString] = currentIndex
				currentString = ""
				currentIndex = index + 1
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

func (st *StringTable) ContainsString(s string) bool {
	_, ok := st.index[s]

	return ok
}

func (st *StringTable) GetIndex(s string) uint32 {
	i, ok := st.index[s]
	if ok {
		return i
	}

	return 0
}

func (st *StringTable) AddString(s string) uint32 {
	if !st.ContainsString(s) {
		var i uint32 = uint32(len(st.content))
		for _, c := range s {
			st.content = append(st.content, byte(c))
		}
		st.content = append(st.content, 0)
		st.index[s] = i
		st.length = uint32(len(st.content))
		return i
	}

	return st.GetIndex(s)
}

func (st *StringTable) Write(writer io.Writer) error {
	err := writeMultiByteUint32(writer, st.length)
	if err == nil {
		if st.length > 0 {
			_, err = writer.Write(st.content)
		}
	}

	return err
}
