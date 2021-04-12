package config

import (
	"net/http"
	"time"
)

// TODO save and restore configs
// TODO webui component for config
// ConfigType holds flags from command line, see flag options for details
type ConfigType struct {
	// UserAgent to use for http client
	UserAgent string
	// FailonError tells worker to exit on any network issues
	FailonError bool
	// KeepTmpFiles keeps intermediate html files after successful mobi creation
	KeepTmpFiles bool
	// LinksFile contains http links in cli mode TODO remove it
	LinksFile string
	// SleepSec seconds to sleep between each http request
	SleepSec int
}

var (
	Config = ConfigType{
		UserAgent: "Mozilla/5.0 (X11; rv:84.0) Gecko/20100101 Firefox/84.0",
		SleepSec:  3,
	}
	Tmpdir = "out"
	// TitlesFname is used to keep track of html titles across runs. Titles are used as mobi chapter titles.
	// TODO  mobi chapter order should be same as inputlinks
	TitlesFname = Tmpdir + "/titles.txt"
	Client      = http.Client{
		Timeout: time.Second * 30,
	}
)
