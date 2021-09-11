package worker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/raviraa/htmltoebook/config"
	"github.com/raviraa/htmltoebook/writer"
)

func (w *Worker) WriteBook() error {
	inpfiles, titleUrls, err := parseTitlesFile(w.conf.TitlesFname())
	if err != nil {
		w.logerr("unable to read intermediate file list. ", err.Error())
		return err
	}
	booktype := config.EbookType
	w.loginfo("Writing ebook " + booktype)
	outfname := path.Join(w.conf.DownloadDir, fmt.Sprintf("%s.%s", w.conf.BookTitle, booktype))

	var book writer.BookWriter
	if booktype == "mobi" {
		book = writer.NewMobi(w.conf.BookTitle, outfname)
	} else {
		book = writer.NewEpub(w.conf.BookTitle, outfname)
	}
	if book == nil {
		w.logerr("Unable to create book writer")
		os.Exit(2)
	}

	if len(inpfiles) == 0 {
		err = errors.New("error fetching any of the links")
		w.logerr(err.Error())
		return err
	}
	w.logsuccess(fmt.Sprintf("Writing %d link(s) to ebook", len(inpfiles)))

	for _, fname := range inpfiles {
		w.loginfo("Adding ", titleUrls[fname][0])
		if titleUrls[fname][1] == ADDIMAGE {
			if err := book.AddImage(fname, titleUrls[fname][0]); err != nil {
				w.logerr("error writing book image " + err.Error())
				return err
			}
		} else {
			b, err := ioutil.ReadFile(fname)
			if err != nil {
				err = fmt.Errorf("error reading intermediate saved html file. %w", err)
				w.logerr(err.Error())
				return err
			}
			// m.NewChapter(titles[fname], b)
			book.AddSection(titleUrls[fname][0], string(b))

		}
	}
	book.AddSection("Index", genIndex(inpfiles, titleUrls))

	if err := book.Write(); err != nil {
		w.logerr("Error writing ebook " + err.Error())
	}
	w.logsuccess("Successfully written " + outfname)

	if !w.conf.KeepTmpFiles {
		w.ClearTmpDir()
	}
	return nil
}

// parseTitlesFile returns file names(fnames), [titles, urls] and error
// fnames is slice of parsed html file names
// titles is map[file name] -> [html title, html url]
// TODO change map value slice to struct
func parseTitlesFile(titleFname string) (fnames []string, titles map[string][]string, err error) {
	titles = make(map[string][]string)
	var lines []string
	lines, err = ReadLines(titleFname)
	if err != nil {
		return
	}

	for _, line := range lines {
		if line == "" {
			continue
		}
		// Each line in the file is of the format "file_name\x00html_title\x00url"
		// TODO use fmt.Sscanf
		spl := strings.SplitN(line, "\x00", 3)
		if len(spl) != 3 {
			err = errors.New("Invalid line in titles file " + line)
			return
		}
		fnames = append(fnames, spl[0])
		titles[spl[0]] = []string{spl[1], spl[2]}
	}
	return
}

func genIndex(inpfiles []string, m map[string][]string) string {
	var s strings.Builder
	s.WriteString("<h3>Index Links</h3>\n<pre>\n")
	// for _, v := range m {
	for _, inpfile := range inpfiles {
		v := m[inpfile]
		if v[1] == ADDIMAGE {
			continue
		}
		// s.WriteString(fmt.Sprintf("<div><a href='%s'>%s</a></div>\n", v[1], v[0]))
		s.WriteString(v[1] + "\n")
	}
	s.WriteString("</pre>")
	return s.String()
}
