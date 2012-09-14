package wbxml

type CodeBook struct {
	TagCodePages       map[byte]CodePage
	TagCodePagesByName map[string]CodePage
	AttributeCodePages map[byte]AttributeCodePage
}

func NewCodeBook() *CodeBook {
	var codeBook *CodeBook = new(CodeBook)
	codeBook.TagCodePages = make(map[byte]CodePage)
	codeBook.TagCodePagesByName = make(map[string]CodePage)
	codeBook.AttributeCodePages = make(map[byte]AttributeCodePage)
	return codeBook
}

func (codeBook *CodeBook) AddTagCodePage(codePage CodePage) {
	if !codeBook.HasTagCode(codePage.Code) && !codeBook.HasNameSpace(codePage.Name) {
		codeBook.TagCodePages[codePage.Code] = codePage
		codeBook.TagCodePagesByName[codePage.Name] = codePage
	}
}

func (codeBook *CodeBook) HasTagCode(code byte) bool {
	var ok bool
	_, ok = codeBook.TagCodePages[code]
	return ok
}

func (codeBook *CodeBook) AddAttributeCodePage(codePage AttributeCodePage) {
	if !codeBook.HasAttributeCode(codePage.Code) {
		codeBook.AttributeCodePages[codePage.Code] = codePage
	}
}

func (codeBook *CodeBook) HasAttributeCode(code byte) bool {
	var ok bool
	_, ok = codeBook.AttributeCodePages[code]
	return ok
}

func (codeBook *CodeBook) IsReady() bool {
	return codeBook.HasTagCode(0)
}

func (codeBook *CodeBook) HasNameSpace(ns string) bool {
	var ok bool
	_, ok = codeBook.TagCodePagesByName[ns]
	return ok
}
