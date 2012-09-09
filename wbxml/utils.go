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
