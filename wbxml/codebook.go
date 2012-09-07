package wbxml

type CodeBook struct {
	CodePages       map[byte]CodePage
	CodePagesByName map[string]CodePage
}

func NewCodeBook() *CodeBook {
	var codeBook *CodeBook = new(CodeBook)
	codeBook.CodePages = make(map[byte]CodePage)
	codeBook.CodePagesByName = make(map[string]CodePage)
	return codeBook
}

func (codeBook *CodeBook) HasCode(Code byte) bool {
	var ok bool
	_, ok = codeBook.CodePages[Code]
	return ok
}

func (codeBook *CodeBook) HasName(Name string) bool {
	var ok bool
	_, ok = codeBook.CodePagesByName[Name]
	return ok
}

func (codeBook *CodeBook) AddCodePage(codePage CodePage) {
	if !codeBook.HasCode(codePage.Code) && !codeBook.HasName(codePage.Name) {
		codeBook.CodePages[codePage.Code] = codePage
		codeBook.CodePagesByName[codePage.Name] = codePage
	}
}
