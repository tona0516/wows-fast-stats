package webapi

type Form struct {
	name    string
	content string
	isFile  bool
}

func NewForm(name string, content string, isFile bool) Form {
	return Form{
		name:    name,
		content: content,
		isFile:  isFile,
	}
}
