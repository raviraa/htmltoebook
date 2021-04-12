package config

import (
	"net/http"
	"time"
)

// ConfigType holds flags from command line, see flag options for details
type ConfigType struct {
	// UserAgent to use for http client
	UserAgent string
	// FailonError tells worker to exit on any network issues
	FailonError bool
	// LinksFile contains http links in cli mode TODO move out
	LinksFile string
	// SleepSec seconds to sleep between each http request
	SleepSec int
}

var (
	Config = ConfigType{
		UserAgent: "Mozilla/5.0 (X11; rv:84.0) Gecko/20100101 Firefox/84.0",
		SleepSec:  33,
	}
	Tmpdir      = "../tmphtmlout"
	TitlesFname = Tmpdir + "/titles.txt"
	Client      = http.Client{
		Timeout: time.Second * 30,
	}
)
