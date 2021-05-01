package writer

type BookWriter interface {
	AddSection(title, body string)
	Write() error
}
