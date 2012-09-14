package wbxml

import (
	"io"
)

func getByteValue(b byte) (byteValue byte, hasContinuationFlag bool) {
	byteValue = b
	hasContinuationFlag = (b&0x80 != 0)
	if hasContinuationFlag {
		byteValue = b &^ 0x80
	}

	return byteValue, hasContinuationFlag
}

func readMultiByteUint32(reader io.ByteReader) (uint32, error) {
	var (
		result       uint32 = 0
		nextByte     byte
		continueRead bool = true
		err          error
	)

	for continueRead {
		nextByte, err = reader.ReadByte()
		if err != nil {
			result = 0
			break
		} else {
			nextByte, continueRead = getByteValue(nextByte)
			result <<= 7
			result |= uint32(nextByte)
		}
	}

	return result, err
}

func writeMultiByteUint32(writer io.Writer, value uint32) error {
	var result []byte

	if value >= 0x80 {
		var temp []byte = make([]byte, 0)
		temp = append(temp, byte(value&0x7F))
		value >>= 7
		for value >= 0x80 {
			temp = append(temp, byte((value&0x7F)|0x80))
			value >>= 7
		}
		temp = append(temp, byte(value|0x80))

		result = make([]byte, 0)
		for i := len(temp) - 1; i >= 0; i-- {
			result = append(result, temp[i])
		}
	} else {
		result = make([]byte, 1)
		result[0] = byte(value)
	}

	_, err := writer.Write(result)
	return err
}
