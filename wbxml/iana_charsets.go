package wbxml

import (
	"fmt"
)

func GetCharsetStringByCode(Code uint32) (string, error) {
	var (
		ok       bool
		charsets map[uint32]string
		result   string
		err      error
	)
	charsets = make(map[uint32]string)

	charsets[3] = "us-ascii"
	charsets[4] = "iso-8859-1"
	charsets[5] = "iso-8859-2"
	charsets[6] = "iso-8859-3"
	charsets[7] = "iso-8859-4"
	charsets[8] = "iso-8859-5"
	charsets[9] = "iso-8859-6"
	charsets[10] = "iso-8859-7"
	charsets[11] = "iso-8859-8"
	charsets[12] = "iso-8859-9"
	charsets[13] = "iso-8859-10"
	charsets[13] = "iso-ir-157"
	charsets[14] = "iso-ir-142"
	charsets[18] = "euc-jp"
	charsets[37] = "iso-2022-kr"
	charsets[38] = "euc-kr"
	charsets[39] = "iso-2022-jp"
	charsets[40] = "iso-2022-jp-2"
	charsets[81] = "iso-8859-6-e"
	charsets[82] = "iso-8859-6-i"
	charsets[84] = "iso-8859-8-e"
	charsets[85] = "iso-8859-8-i"
	charsets[106] = "utf-8"
	charsets[109] = "iso-8859-13"
	charsets[110] = "iso-8859-14"
	charsets[111] = "iso-8859-15"
	charsets[112] = "iso-8859-16"
	charsets[1012] = "utf-7"
	charsets[1013] = "utf-16be"
	charsets[1014] = "utf-16le"
	charsets[1015] = "utf-16"
	charsets[1017] = "utf-32"
	charsets[1018] = "utf-32be"
	charsets[1019] = "utf-32le"
	charsets[2025] = "gb2312"
	charsets[2026] = "big5"
	charsets[2084] = "koi8-r"

	result, ok = charsets[Code]
	if !ok {
		result = ""
		err = fmt.Errorf("Charset code %d not found!", Code)
	}

	return result, err
}
