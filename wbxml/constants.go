package wbxml

const (
	SWITCH_PAGE        byte = 0x00
	END                byte = 0x01
	ENTITY             byte = 0x02
	STR_I              byte = 0x03
	LITERAL            byte = 0x04
	EXT_I_0            byte = 0x40
	EXT_I_1            byte = 0x41
	EXT_I_2            byte = 0x42
	PI                 byte = 0x43
	LITERAL_C          byte = 0x44
	EXT_T_0            byte = 0x80
	EXT_T_1            byte = 0x81
	EXT_T_2            byte = 0x82
	STR_T              byte = 0x83
	LITERAL_A          byte = 0x84
	EXT_0              byte = 0xC0
	EXT_1              byte = 0xC1
	EXT_2              byte = 0xC2
	OPAQUE             byte = 0xC3
	LITERAL_AC         byte = 0xC4
	WBXML_1_3          byte = 0x03
	UNKNOWN_PI         byte = 0x01
	CHARSET_UTF8       byte = 0x6A
	CHARSET_UNKNOWN    byte = 0x00
	TAG_HAS_ATTRIBUTES byte = 0x80
	TAG_HAS_CONTENT    byte = 0x40
)
