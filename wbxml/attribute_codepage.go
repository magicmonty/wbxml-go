package wbxml

import (
	"fmt"
	"strings"
)

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
	ValueCodes map[string]byte
}

func NewAttributeCodePage(code byte) AttributeCodePage {
	var codePage AttributeCodePage
	codePage.Code = code
	codePage.Attributes = make(map[byte]AttributeCodePageEntry)
	codePage.Values = make(map[byte]string)
	codePage.ValueCodes = make(map[string]byte)
	return codePage
}

func (codePage *AttributeCodePage) AddAttribute(name string, valuePrefix string, code byte) {
	if !codePage.HasCode(code) && code < ATTRIBUTE_VALUE_SPACE_START {
		var entry AttributeCodePageEntry
		entry.Name = name
		entry.ValuePrefix = valuePrefix
		codePage.Attributes[code] = entry
	}
}

func (codePage *AttributeCodePage) AddAttributeValue(value string, code byte) {
	if !codePage.HasValueCode(code) && !codePage.HasValue(value) && code >= ATTRIBUTE_VALUE_SPACE_START {
		codePage.Values[code] = value
		codePage.ValueCodes[value] = code
	}
}

func (codePage *AttributeCodePage) HasCode(code byte) bool {
	var ok bool
	_, ok = codePage.Attributes[code]

	return ok
}

func (codePage *AttributeCodePage) HasValueCode(code byte) bool {
	var ok bool
	_, ok = codePage.Values[code]

	return ok
}

func (codePage *AttributeCodePage) HasValue(value string) bool {
	var ok bool
	_, ok = codePage.ValueCodes[value]

	return ok
}

func (codePage *AttributeCodePage) GetString(code byte) string {
	var (
		result string = ""
		entry  AttributeCodePageEntry
	)

	if code < ATTRIBUTE_VALUE_SPACE_START {
		if codePage.HasCode(code) {
			entry = codePage.Attributes[code]
			result = " " + entry.Name + "=\""
			if entry.ValuePrefix != "" {
				result += entry.ValuePrefix
			}
		}
	} else {
		if codePage.HasValueCode(code) {
			result = codePage.Values[code]
		}
	}

	return result
}

func (codePage *AttributeCodePage) HasAttribute(name string, value string) bool {
	for _, a := range codePage.Attributes {
		if a.Name == name {
			if value == "" {
				if a.ValuePrefix == "" {
					return true
				}
			} else {
				if a.ValuePrefix == "" || strings.HasPrefix(value, a.ValuePrefix) {
					return true
				}
			}
		}
	}

	return false
}

func (codePage *AttributeCodePage) Tokenize(name string, value string) (byte, []string, error) {
	var (
		attributeToken byte     = 0x00
		valueTokens    []string = nil
		err            error    = nil
	)

	if codePage.HasAttribute(name, value) {
		attributeToken = codePage.getAttributeId(name, value)
		valueTokens = make([]string, 0)
		if strings.HasPrefix(value, codePage.Attributes[attributeToken].ValuePrefix) {
			value = strings.Replace(value, codePage.Attributes[attributeToken].ValuePrefix, "", -1)
		}

		splittedValues := codePage.splitValues(value)
		for _, sv := range splittedValues {
			valueTokens = append(valueTokens, sv)
		}
	} else {
		err = fmt.Errorf("codepage has no matching attribute entry (name='%s', value='%s')", name, value)
	}

	return attributeToken, valueTokens, err
}

func (codePage *AttributeCodePage) getAttributeId(name string, value string) byte {
	if value == "" {
		for attributeId, a := range codePage.Attributes {
			if a.Name == name && a.ValuePrefix == "" {
				return attributeId
			}
		}
	} else {
		for attributeId, a := range codePage.Attributes {
			if a.Name == name && a.ValuePrefix != "" && strings.HasPrefix(value, a.ValuePrefix) {
				return attributeId
			}
		}

		for attributeId, a := range codePage.Attributes {
			if a.Name == name && a.ValuePrefix == "" {
				return attributeId
			}
		}
	}
	return 0
}

func (codePage *AttributeCodePage) splitValues(value string) []string {
	var result []string = make([]string, 0)
	if value == "" {
		return result
	}

	for _, v := range codePage.Values {
		if strings.Contains(value, v) {
			splittedStrings := strings.Split(value, v)
			for index, splittedString := range splittedStrings {
				if splittedString != "" {
					for _, sv := range codePage.splitValues(splittedString) {
						if sv != "" {
							result = append(result, sv)
						}
					}
					if index < len(splittedStrings)-1 {
						result = append(result, v)
					}
				}
			}
			break
		}
	}

	if len(result) == 0 {
		result = append(result, value)
	}

	return result
}
