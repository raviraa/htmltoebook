package web

import (
	"context"
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jfyne/live"
	"github.com/raviraa/htmltoebook/config"
	"github.com/raviraa/htmltoebook/types"
	"github.com/raviraa/htmltoebook/worker"
	opengolang "github.com/skratchdot/open-golang/open"
)

// state shared across ui clients.
var workerRunning bool
var runningMu sync.Mutex

type model struct {
	LogMsgs  []types.LogMsg
	ShowConf bool
	Conf     *config.ConfigType
	Running  bool

	cancel context.CancelFunc
	worker *worker.Worker
}

var liveHandler *live.Handler

//go:embed root.html
var rootHtml string

//go:embed milligram.css
var cssFramework []byte

//go:embed favicon.png
var favIconData []byte

//go:embed progress.gif
var progressGif []byte

func newModel(s *live.Socket) *model {
	m, ok := s.Assigns().(*model)
	// on first socket connect when client loads page
	if !ok {
		conf := config.New()
		return &model{
			worker: worker.New(liveHandler, conf),
			Conf:   conf,
		}
	}
	return m
}

func NewWeb() {
	t, err := template.New("root.html").Parse(rootHtml)
	if err != nil {
		log.Fatal(err)
	}
	// App should run only on localhost, security for public network not considered
	h, err := live.NewHandler(live.NewCookieStore("session-name", []byte("nosecret")), live.WithTemplateRenderer(t))
	if err != nil {
		log.Fatal(err)
	}
	liveHandler = h
	setEvents(h)

	http.Handle("/", h)
	mystatic := webstatic{}
	http.Handle("/favicon.ico", mystatic)
	http.Handle("/milligram.css", mystatic)
	http.Handle("/progress.gif", mystatic)
	http.Handle("/quit", mystatic)
	http.Handle("/live.js", live.Javascript{})
	http.Handle("/auto.js.map", live.JavascriptMap{})
	port := os.Getenv("PORT")
	if port == "" {
		port = "localhost:39394"
	}
	log.Println("listening on ", port)
	go func() {
		time.Sleep(2 * time.Second)
		// launch link in browser
		localhref := "http://" + port
		if err = opengolang.Run(localhref); err != nil {
			log.Println("Error opening app in browser. Please open the link manually: ", localhref)
		}
	}()
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

type webstatic struct{}

func (webstatic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.RequestURI {
	case "/favicon.ico":
		w.Header().Set("Content-Type", "image/png")
		w.Write(favIconData)
	case "/milligram.css":
		w.Header().Set("Content-Type", "text/css")
		w.Write(cssFramework)
	case "/progress.gif":
		w.Header().Set("Content-Type", "image/gif")
		w.Write(progressGif)
	case "/quit":
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h3>Please close the tab</h3>"))
		go quitApplication()
	}
}

func quitApplication() {
	time.Sleep(3 * time.Second)
	log.Println("Exiting application")
	os.Exit(0)
}
