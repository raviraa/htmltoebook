package web

import (
	"context"
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
)

type model struct {
	LogMsgs    []types.LogMsg
	Running    bool
	InputLinks []string
}

func newModel(s *live.Socket) *model {
	m, ok := s.Assigns().(*model)
	if !ok {
		return &model{}
	}
	return m
}

func NewWeb() {
	t, err := template.ParseFiles("root.html", "view.html")
	if err != nil {
		log.Fatal(err)
	}

	h, err := live.NewHandler(live.NewCookieStore("session-name", []byte("weak-secret")), live.WithTemplateRenderer(t))
	if err != nil {
		log.Fatal(err)
	}
	worker.SetHandler(h)

	h.Mount = onMount
	// h.HandleEvent(evstart, onLogMsg)
	h.HandleSelf(evlogmsg, onLogMsg)

	http.Handle("/", h)
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	port := os.Getenv("PORT")
	if port == "" {
		port = "localhost:9394"
	}
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func onMount(ctx context.Context, r *http.Request, s *live.Socket) (interface{}, error) {
	m := newModel(s)
	if s.Connected() {
		worker.StartWorker()
	}
	return m, nil
}
