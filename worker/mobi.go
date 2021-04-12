package worker

import (
	"io/ioutil"
	"localhost/htmltoebook/config"
	"os"
	"strings"

	"github.com/766b/mobi"
)

func WriteMobi() {
	inpfiles, titles := parseTitlesFile()
	loginfo("Writing mobi file")
	outfname := config.Tmpdir + "/output.mobi"
	m, err := mobi.NewWriter(outfname)
	if err != nil {
		panic(err)
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

	for _, fname := range inpfiles {
		loginfo("Adding ", titles[fname])
		b, err := ioutil.ReadFile(fname)
		panicerr(err)
		m.NewChapter(titles[fname], b)
	}
	// Output MOBI File
	m.Write()
	logsuccess("Sucessfully written " + outfname)

	if !config.Config.KeepTmpFiles {
		loginfo("Cleaning temporary files")
		for _, fname := range inpfiles {
			os.Remove(fname)
		}
		os.Remove(config.TitlesFname)

	}
}

func parseTitlesFile() ([]string, map[string]string) {
	titles := make(map[string]string)
	fnames := make([]string, 0)
	lines := ReadLines(config.TitlesFname)
	for _, line := range lines {
		if line == "" {
			continue
		}
		spl := strings.SplitN(line, " ", 2)
		if len(spl) != 2 {
			panic("Invalid line in titles file " + line)
		}
		fnames = append(fnames, spl[0])
		titles[spl[0]] = spl[1]
	}
	return fnames, titles
}
