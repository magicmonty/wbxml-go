package wbxml

type CodeBook struct {
	TagCodePages       map[byte]CodePage
	AttributeCodePages map[byte]AttributeCodePage
}

func NewCodeBook() *CodeBook {
	var codeBook *CodeBook = new(CodeBook)
	codeBook.TagCodePages = make(map[byte]CodePage)
	codeBook.AttributeCodePages = make(map[byte]AttributeCodePage)
	return codeBook
}

func (codeBook *CodeBook) HasTagCode(Code byte) bool {
	var ok bool
	_, ok = codeBook.TagCodePages[Code]
	return ok
}

func (codeBook *CodeBook) HasAttributeCode(Code byte) bool {
	var ok bool
	_, ok = codeBook.AttributeCodePages[Code]
	return ok
}

func (codeBook *CodeBook) AddTagCodePage(codePage CodePage) {
	if !codeBook.HasTagCode(codePage.Code) {
		codeBook.TagCodePages[codePage.Code] = codePage
	}
}

func (codeBook *CodeBook) AddAttributeCodePage(codePage AttributeCodePage) {
	if !codeBook.HasAttributeCode(codePage.Code) {
		codeBook.AttributeCodePages[codePage.Code] = codePage
	}
}

func (codeBook *CodeBook) IsReady() bool {
	return codeBook.HasTagCode(0)
}
