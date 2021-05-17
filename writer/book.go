package writer

type BookWriter interface {
	AddSection(title, body string)
	// imagefile is location of the image on disk
	// imagehtmsrc is src attribute of <img> set in source html
	AddImage(imagefile, imageattrsrc string) error
	Write() error
}
