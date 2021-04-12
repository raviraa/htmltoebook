package web

import (
	"html/template"
	"localhost/htmltoebook/types"
	"localhost/htmltoebook/worker"
	"log"
	"net/http"
	"os"

	"github.com/jfyne/live"
)

const (
	evlogmsg = types.Evlogmsg
	evstart  = "evstart"
	evclear  = "evclear"
	evstop   = "evstop"
)

type model struct {
	LogMsgs []types.LogMsg
	Running bool
	// InputLinks string
}

func newModel(s *live.Socket) *model {
	m, ok := s.Assigns().(*model)
	if !ok {
		return &model{}
	}
	return m
}

func NewWeb() {
	// TODO embed template https://golang.org/pkg/embed/
	t, err := template.ParseFiles("root.html")
	if err != nil {
		log.Fatal(err)
	}

	h, err := live.NewHandler(live.NewCookieStore("session-name", []byte("weak-secret")), live.WithTemplateRenderer(t))
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
	if err != nil {
		log.Fatal(err)
	}
}
