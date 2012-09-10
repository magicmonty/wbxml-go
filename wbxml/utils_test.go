package wbxml

import (
	"fmt"
)

func ExampleReadMultiByteUint32() {
	var (
		result uint32
	)

	result, _ = readMultiByteUint32(MakeDataBuffer(0x81, 0x20))
	fmt.Printf("%d\n", result)
	result, _ = readMultiByteUint32(MakeDataBuffer(0x60))
	fmt.Printf("%d\n", result)
	// OUTPUT: 160
	// 96
}
