package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type ConfigType struct {
	UserAgent string `comment:"UserAgent to use for http client"`

	FailonError bool `comment:"Exit immediately on any network errors"`

	KeepTmpFiles bool `comment:"Do not remove intermediate html files after successful ebook creation"`

	SleepSec int `comment:"SleepSec seconds to sleep between each http request"`

	Tmpdir string `comment:"Directory to keep downloaded web pages and generated ebook"`

	BookTitle string `comment:"Title to be used in the ebook"`

	PreBreaks bool `comment:"Add <br> for each line in <pre> block. Needed for some old ebook readers"`

	IncludeImages bool `comment:"Download images included in web page and add them to book"`
}

func New() *ConfigType {
	homedir, _ := os.UserHomeDir()
	config := ConfigType{
		// defaults when config file is absent
		UserAgent: "Mozilla/5.0",
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
	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	return enc.Encode(c)
}

func (c *ConfigType) readConf() error {
	// _, err := toml.DecodeFile(confLocation(), c)
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
	return c.Tmpdir + "/titles.txt"
}
