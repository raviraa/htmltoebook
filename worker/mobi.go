package worker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/766b/mobi"
)

func (w *Worker) WriteMobi() error {
	inpfiles, titles, err := parseTitlesFile(w.conf.TitlesFname())
	if err != nil {
		w.logerr("unable to read intermediate file list. ", err.Error())
		return err
	}
	w.loginfo("Writing mobi file")
	outfname := w.conf.Tmpdir + "/output.mobi"
	m, err := mobi.NewWriter(outfname)
	if err != nil {
		w.logerr("Failed opening mobi file ", outfname, err.Error())
		return err
	}

	// TODO add title, image to configuration
	m.Title("Book Title")
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

// parseTitlesFile Returns fnames, titles and error
// fnames is slice of parsed html file names
// titles is map[file name] -> html title
func parseTitlesFile(titleFname string) (fnames []string, titles map[string]string, err error) {
	titles = make(map[string]string)
	var lines []string
	lines, err = ReadLines(titleFname)
	if err != nil {
		return
	}

	for _, line := range lines {
		if line == "" {
			continue
		}
		spl := strings.SplitN(line, "\x00", 2)
		if len(spl) != 2 {
			err = errors.New("Invalid line in titles file " + line)
			return
		}
		fnames = append(fnames, spl[0])
		titles[spl[0]] = spl[1]
	}
	return
}
