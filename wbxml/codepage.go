package wbxml

import (
	"fmt"
)

type CodePage struct {
	Name     string
	Code     byte
	Tags     map[byte]string
	TagCodes map[string]byte
}

func NewCodePage(Name string, Code byte) CodePage {
	var codePage CodePage
	codePage.Name = Name
	codePage.Code = Code
	codePage.Tags = make(map[byte]string)
	codePage.TagCodes = make(map[string]byte)
	return codePage
}

func (codePage *CodePage) AddTag(Tag string, Code byte) {
	if !codePage.HasTagCode(Code) && !codePage.HasTag(Tag) {
		codePage.Tags[Code] = Tag
		codePage.TagCodes[Tag] = Code
	}
}

func (codePage *CodePage) HasTag(Tag string) bool {
	var ok bool
	_, ok = codePage.TagCodes[Tag]

	return ok
}

func (codePage *CodePage) HasTagCode(Code byte) bool {
	var ok bool
	_, ok = codePage.Tags[Code]

	return ok
}

func (codePage *CodePage) GetNameSpaceString() string {
	var (
		part1 string = ""
		part2 string = ""
	)

	if codePage.Code > 0 {
		if codePage.Code/26 > 0 {
			part1 = fmt.Sprintf("%c", 0x41+((codePage.Code/26)-1)%26)
		}

		part2 = fmt.Sprintf("%c", 0x41+(codePage.Code%26))
	}

	return part1 + part2
}

func (codePage *CodePage) GetNameSpaceDeclaration() string {
	var (
		nameSpacePrefix string
		result          string
	)

	if codePage.Name != "" {
		nameSpacePrefix = codePage.GetNameSpaceString()
		if nameSpacePrefix == "" {
			result = " xmlns=\""
		} else {
			result = " xmlns:" + nameSpacePrefix + "=\""
		}

		result += codePage.Name + "\""
	}

	return result
}
