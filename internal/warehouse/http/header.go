package http

type Header struct {
	string
}

func NewHeader(name, value string) string {
	endOfHeader := "\r\n"
	header := name + ":" + value + endOfHeader
	return header
}
