package reader

type BookReader interface {
  ReadFiles(func(sectionData []byte)) error
}
