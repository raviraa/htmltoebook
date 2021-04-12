package web

import (
	"context"
	_ "embed"
	"html/template"
	"localhost/htmltoebook/types"
	"localhost/htmltoebook/worker"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jfyne/live"
)

// Events to handle in from ui or worker
const (
	evlogmsg        = types.Evlogmsg
	evWorkerStopped = types.EvWorkerStopped
	evstart         = "evstart"
	evclear         = "evclear"
	evstop          = "evstop"
)

// state shared across ui clients
var workerRunning bool

type model struct {
	LogMsgs []types.LogMsg
	cancel  context.CancelFunc
}

func newModel(s *live.Socket) *model {
	m, ok := s.Assigns().(*model)
	if !ok {
		return &model{}
	}
	return m
}

//go:embed root.html
var rootHtml string

//go:embed milligram.css
var cssFramework string

func NewWeb() {
	// Using embedded root.html and css framework  in binary
	rootHtml = strings.Replace(rootHtml, "/*CSSFRAMEWORK*/", cssFramework, 1)
	t, err := template.New("root.html").Parse(rootHtml)
	if err != nil {
		log.Fatal(err)
	}
	hostname, _ := os.Hostname()
	// App should run only on localhost, security for public network not considered
	h, err := live.NewHandler(live.NewCookieStore("session-name", []byte("secret"+hostname)), live.WithTemplateRenderer(t))
	if err != nil {
		log.Fatal(err)
	}
	worker.SetHandler(h)
	setEvents(h)

	http.Handle("/", h)
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	port := os.Getenv("PORT")
	if port == "" {
		port = "localhost:9394"
	}
	log.Println("listening on ", port)
	err = http.ListenAndServe(port, nil)
	// TODO launch link in browser
	if err != nil {
		log.Fatal(err)
	}
}
