package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	// "github.com/pelletier/go-toml"
	"encoding/json"
)

var (
	// IdleTimeout is timeout for chrome page download
	IdleTimeout = time.Second * 9

	// RetryNum is number of times to try each failed link
	RetryNum = 3

	// ChromeRestartNum is number of requests to fetch before restarting chrome
	ChromeRestartNum = 9

	//Ebook type, epub or mobi.
	EbookType = "epub"
)

type ConfigType struct {
	UserAgent string `comment:"UserAgent to use for http client"`

	FailonError bool `comment:"Exit immediately on any network errors"`

	// Do not remove intermediate html files after successful ebook creation
	KeepTmpFiles bool `toml:"-"`

	SleepSec int `comment:"SleepSec seconds to sleep between each http request"`

	DownloadDir string `comment:"Directory to keep downloaded web pages and generated ebook"`

	BookTitle string `comment:"Title to be used in the ebook"`
	PreBreaks bool   `comment:"Add <br> for each line in <pre> block. Needed for some old ebook readers"`

	IncludeImages bool `comment:"Download images included in web page and add them to book"`

	ChromeDownload bool `comment:"Download web pages using chrome browser"`

	UseTextParser bool `comment:"Use text content from html parser. Useful for epub books with parsing issues"`

	MinPageSize int `comment:"If size of page content is less than specified size, program exits"`

	Links string `toml:"-"`
}

func New() *ConfigType {
	homedir, _ := os.UserHomeDir()
	config := ConfigType{
		// defaults when config file is absent
		// cmd editor read/write discards these values. set defaults on WriteConf
		UserAgent:   "Mozilla/5.0",
		SleepSec:    9,
		MinPageSize: 99,
		BookTitle:   "Book Title",
		DownloadDir: path.Join(homedir, "Downloads", "htmltoebook"),
	}
	config.readConf()
	return &config
}

// Write json config to file.
// #Links is written in json file, but not user editable toml file. toml format is used for cmd/editor logic
func (c *ConfigType) WriteConf() error {
	f, err := os.OpenFile(confLocation(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	// enc := toml.NewEncoder(f)
	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")

	return enc.Encode(*c)
}

func (c *ConfigType) Tmpdir() string {
	// return path.Join(c.DownloadDir, c.BookTitle)
	return path.Join(os.TempDir(), c.BookTitle)
}

func (c *ConfigType) readConf() error {
	log.Println("Using config ", confLocation())
	b, err := ioutil.ReadFile(confLocation())
	if err != nil {
		log.Println("Failure in reading config ", err)
		return err
	}
	if err = json.Unmarshal(b, c); err != nil {
		log.Println("Failure in parsing config ", err)
		return err
	}
	return nil
}

const confFileName = ".htmltoebook.json"

func confLocation() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return confFileName
	}
	return path.Join(home, confFileName)
}

// TitlesFname file is used to keep track of stripped html files and titles across runs.
// Titles are used as book chapter titles.
// Each line in the file is of the format "file_name\x00html_title\x00url"
func (c *ConfigType) TitlesFname() string {
	// return c.Tmpdir() + "/titles.txt"
	return path.Join(c.Tmpdir(), "titles.txt")
}
