package writer

/*
import  "github.com/766b/mobi"
TODO
func (w *Worker) WriteMobi() error {
	inpfiles, titles, err := parseTitlesFile(w.conf.TitlesFname())
	if err != nil {
		w.logerr("unable to read intermediate file list. ", err.Error())
		return err
	}
	w.loginfo("Writing mobi file")
	outfname := fmt.Sprintf("%s/%s.mobi", w.conf.Tmpdir, w.conf.BookTitle)
	m, err := mobi.NewWriter(outfname)
	if err != nil {
		w.logerr("Failed opening mobi file ", outfname, err.Error())
		return err
	}

	m.Title(w.conf.BookTitle)
	m.Compression(mobi.CompressionNone) // LZ77 compression is also possible using  mobi.CompressionPalmDoc
	// m.Compression(mobi.CompressionPalmDoc)

	// Add cover image
	// m.AddCover("data/cover.jpg", "data/thumbnail.jpg")

	// Meta data
	m.NewExthRecord(mobi.EXTH_DOCTYPE, "EBOK")
	m.NewExthRecord(mobi.EXTH_AUTHOR, "Book Author Name")

	if len(inpfiles) == 0 {
		err = errors.New("error fetching any of the links")
		w.logerr(err.Error())
		return err
	}
	w.logsuccess(fmt.Sprintf("Writing %d link(s) to mobi file", len(inpfiles)))

	for _, fname := range inpfiles {
		w.loginfo("Adding ", titles[fname])
		b, err := ioutil.ReadFile(fname)
		if err != nil {
			err = fmt.Errorf("error reading intermediate saved html file. %w", err)
			w.logerr(err.Error())
			return err
		}
		m.NewChapter(titles[fname], b)
	}
	// Output MOBI File
	m.Write()
	w.logsuccess("Sucessfully written " + outfname)

	if !w.conf.KeepTmpFiles {
		w.loginfo("Cleaning temporary files")
		for _, fname := range inpfiles {
			os.Remove(fname)
		}
		os.Remove(w.conf.TitlesFname())

	}
	return nil
}


*/
