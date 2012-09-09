package wbxml

const (
	ATTRIBUTE_VALUE_SPACE_START byte = 0x80
)

type AttributeCodePageEntry struct {
	Name        string
	ValuePrefix string
}

type AttributeCodePage struct {
	Code       byte
	Attributes map[byte]AttributeCodePageEntry
	Values     map[byte]string
}

func NewAttributeCodePage(Code byte) AttributeCodePage {
	var codePage AttributeCodePage
	codePage.Code = Code
	codePage.Attributes = make(map[byte]AttributeCodePageEntry)
	codePage.Values = make(map[byte]string)
	return codePage
}

func (codePage *AttributeCodePage) AddAttribute(Name string, ValuePrefix string, Code byte) {
	if !codePage.HasCode(Code) && Code < ATTRIBUTE_VALUE_SPACE_START {
		var entry AttributeCodePageEntry
		entry.Name = Name
		entry.ValuePrefix = ValuePrefix
		codePage.Attributes[Code] = entry
	}
}

func (codePage *AttributeCodePage) AddAttributeValue(Value string, Code byte) {
	if !codePage.HasValueCode(Code) && Code >= ATTRIBUTE_VALUE_SPACE_START {
		codePage.Values[Code] = Value
	}
}

func (codePage *AttributeCodePage) HasCode(Code byte) bool {
	var ok bool
	_, ok = codePage.Attributes[Code]

	return ok
}

func (codePage *AttributeCodePage) HasValueCode(Code byte) bool {
	var ok bool
	_, ok = codePage.Values[Code]

	return ok
}

func (codePage *AttributeCodePage) GetString(Code byte) string {
	var (
		result string = ""
		entry  AttributeCodePageEntry
	)

	if Code < ATTRIBUTE_VALUE_SPACE_START {
		if codePage.HasCode(Code) {
			entry = codePage.Attributes[Code]
			result = " " + entry.Name + "="
			if entry.ValuePrefix != "" {
				result += "\"" + entry.ValuePrefix
			}
		}
	} else {
		if codePage.HasValueCode(Code) {
			result = codePage.Values[Code]
		}
	}

	return result
}
