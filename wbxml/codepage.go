package wbxml

type CodePage struct {
	Name  string
	Code  byte
	Codes map[byte]string
	Tags  map[string]byte
}

func NewCodePage(Name string, Code byte) CodePage {
	var codePage CodePage
	codePage.Name = Name
	codePage.Code = Code
	codePage.Codes = make(map[byte]string)
	codePage.Tags = make(map[string]byte)
	return codePage
}

func (codePage *CodePage) AddItem(Tag string, Code byte) {
	if !codePage.HasCode(Code) && !codePage.HasTag(Tag) {
		codePage.Codes[Code] = Tag
		codePage.Tags[Tag] = Code
	}
}

func (codePage *CodePage) HasTag(Tag string) bool {
	var ok bool
	_, ok = codePage.Tags[Tag]

	return ok
}

func (codePage *CodePage) HasCode(Code byte) bool {
	var ok bool
	_, ok = codePage.Codes[Code]

	return ok
}
