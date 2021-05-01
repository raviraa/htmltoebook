package config

import (
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

// ConfigType holds flags from command line, see flag options for details
type ConfigType struct {
	// UserAgent to use for http client
	UserAgent string
	// FailonError tells worker to exit on any network issues
	FailonError bool
	// KeepTmpFiles keeps intermediate html files after successful mobi creation
	KeepTmpFiles bool
	// SleepSec seconds to sleep between each http request
	SleepSec int
	// Directory to keep downloaded web pages and generated ebook
	Tmpdir string
	// Title to be used in the ebook
	BookTitle string
	// add <br> for each line in <pre> block
	AddPreBreaks bool
}

func New() *ConfigType {
	homedir, _ := os.UserHomeDir()
	config := ConfigType{
		// defaults when config file is absent
		UserAgent: "Mozilla/5.0 (X11; rv:84.0) Gecko/20100101 Firefox/84.0",
		SleepSec:  3,
		BookTitle: "Book Title",
		Tmpdir:    path.Join(homedir, "Downloads", "htmltoebook"),
	}
	config.readConf()
	return &config
}

func (c *ConfigType) WriteConf() error {
	f, err := os.OpenFile(confLocation(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	enc := toml.NewEncoder(f)
	return enc.Encode(c)
}

func (c *ConfigType) readConf() error {
	_, err := toml.DecodeFile(confLocation(), c)
	if err != nil {
		log.Println("Failure in reading config ", err)
		return err
	}
	return nil
}

const confFileName = ".htmltoebook.toml"

func confLocation() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return confFileName
	}
	return path.Join(home, confFileName)
}

// TitlesFname file is used to keep track of stripped html files and titles across runs.
// Titles are used as mobi chapter titles.
// Each line in the file is of the format "file_name\x00html_title\x00url"
// TODO  mobi chapter order should be same as inputlinks, problem with failed runs
func (c *ConfigType) TitlesFname() string {
	return c.Tmpdir + "/titles.txt"
}
