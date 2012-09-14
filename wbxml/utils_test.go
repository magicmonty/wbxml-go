package wbxml

import (
	"bytes"
	"fmt"
	"testing"
)

func ExampleReadMultiByteUint32() {
	var (
		result uint32
	)

	result, _ = readMultiByteUint32(makeDataBuffer(0x81, 0x20))
	fmt.Printf("%d\n", result)
	result, _ = readMultiByteUint32(makeDataBuffer(0x60))
	fmt.Printf("%d\n", result)
	result, _ = readMultiByteUint32(makeDataBuffer(0x81, 0x00))
	fmt.Printf("%d\n", result)
	// OUTPUT: 160
	// 96
	// 128
}

func Test_ReadMultiByteUint32_Failure(t *testing.T) {
	var (
		result uint32
		err    error
	)

	result, err = readMultiByteUint32(makeDataBuffer(0x81))

	if err == nil {
		t.Error("Error should be set but was nil")
	}

	if result != 0 {
		t.Errorf("Result should be 0 but was %d", result)
	}
}

func Example_WriteMulitByteUint32() {
	b := bytes.NewBuffer([]byte{})
	writeMultiByteUint32(b, 0x60)
	writeMultiByteUint32(b, 0x7F)
	writeMultiByteUint32(b, 0x80)
	writeMultiByteUint32(b, 160)

	printByteStream(b)
	// OUTPUT: 60 7F 81 00 81 20
}
