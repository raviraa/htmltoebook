package main

import (
	"net/http"
	"time"
)

// Config holds flags from command line, see flag options for details
type Config struct {
	UserAgent     string
	FailonError   bool
	LinksFile     string
	SleepInterval int
}

var (
	config      = Config{}
	tmpdir      = "out"
	titlesfname = tmpdir + "/titles.txt"
	client      = http.Client{
		Timeout: time.Second * 30,
	}
)
