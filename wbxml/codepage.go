package wbxml

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
