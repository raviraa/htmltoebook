package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/raviraa/htmltoebook/htmlextract"
	"github.com/raviraa/htmltoebook/types"
	"github.com/raviraa/htmltoebook/worker"

	"github.com/jfyne/live"
)

// Events to handle in from ui or worker
const (
	evlogmsg        = types.Evlogmsg
	evWorkerStopped = types.EvWorkerStopped
	// Starts downloading and processing web links
	evstart = "evstart"
	// Clear tmp files and removes tmpdir
	evclear = "evclear"
	// Request to stop processing
	evstop = "evstop"
	// Save Settings form data
	evsave = "evsave"
	// Show/Hide Settings modal
	evconf = "evconf"
	// Show/Hide Html Snippet modal
	evhtmlsnippet = "evhtmlsnippet"
	// Form submit for html snippet
	evsnippetsave = "evsnippetsave"
)

func setEvents(h *live.Handler) {
	h.Mount = onMount
	h.HandleEvent(evstart, onStart)
	h.HandleEvent(evstop, onStop)
	h.HandleEvent(evclear, onClear)
	h.HandleEvent(evsave, onSave)
	h.HandleEvent(evconf, onConf)
	h.HandleEvent(evhtmlsnippet, onHtmlSnippet)
	h.HandleEvent(evsnippetsave, onHtmlSnippetSave)

	h.HandleSelf(evlogmsg, onLogMsg)
	h.HandleSelf(evWorkerStopped, onWorkerStopped)
}

func onMount(ctx context.Context, r *http.Request, s *live.Socket) (interface{}, error) {
	m := newModel(s)
	if s.Connected() {
		go func() {
			log.Println("Connected")
			m.worker.AppendLog("Connected. Ready to start.", "info")

		}()
	}
	return m, nil
}

func onLogMsg(ctx context.Context, s *live.Socket, p map[string]interface{}) (interface{}, error) {
	m := newModel(s)
	msgparam, ok := p["logmsg"]
	if !ok {
		return m, fmt.Errorf("log not provided %v", p)
	}
	logmsg, ok := msgparam.(types.LogMsg)
	if !ok {
		return m, fmt.Errorf("log conversion error %v", p)
	}
	// m.LogMsgs = []types.LogMsg{logmsg}
	m.LogMsgs = append(m.LogMsgs, logmsg)
	return m, nil
}

func onClear(ctx context.Context, s *live.Socket, p map[string]interface{}) (interface{}, error) {
	m := newModel(s)
	m.LogMsgs = nil
	m.worker.ClearTmpDir()
	return m, nil
}

func onConf(ctx context.Context, s *live.Socket, p map[string]interface{}) (interface{}, error) {
	m := newModel(s)
	m.ShowConf = !m.ShowConf
	return m, nil
}

func onHtmlSnippet(ctx context.Context, s *live.Socket, p map[string]interface{}) (interface{}, error) {
	m := newModel(s)
	m.ShowHtmlSnippet = !m.ShowHtmlSnippet
	return m, nil
}

func onStop(ctx context.Context, s *live.Socket, p map[string]interface{}) (interface{}, error) {
	m := newModel(s)
	if m.cancel != nil {
		m.worker.AppendLog("Requesting for cancellation", "warn")
		m.cancel()
	}
	return m, nil
}

func onStart(ctx context.Context, s *live.Socket, p map[string]interface{}) (interface{}, error) {
	m := newModel(s)

	// links textarea name could have random number at the end due to HtmlSnippet
	var linksParam string
	for k, v := range p {
		if k[:5] == "links" {
			linksParam = v.(string)
		}
	}
	links := worker.SplitLinks(linksParam)
	if len(links) > 0 {
		runningMu.Lock()
		if workerRunning {
			m.worker.AppendLog("Worker already running", "warn")
			return m, nil
		}
		workerRunning = true
		m.Running = true
		runningMu.Unlock()

		m.Conf.BookTitle = live.ParamString(p, "BookTitle")
		var ctx context.Context
		ctx, m.cancel = context.WithCancel(context.Background())
		m.worker.StartWorker(ctx, links)
		return m, nil
	}
	m.worker.AppendLog("Please add http links in the text area to process", "info")
	return m, nil
}

func onWorkerStopped(ctx context.Context, s *live.Socket, p map[string]interface{}) (interface{}, error) {
	m := newModel(s)
	runningMu.Lock()
	workerRunning = false
	m.Running = false
	runningMu.Unlock()
	log.Println("worked stopped from ui")
	return m, nil
}

func onSave(ctx context.Context, s *live.Socket, p map[string]interface{}) (interface{}, error) {
	m := newModel(s)
	m.Conf.SleepSec = live.ParamInt(p, "SleepSec")
	m.Conf.FailonError = live.ParamCheckbox(p, "FailonError")
	m.Conf.KeepTmpFiles = live.ParamCheckbox(p, "KeepTmpFiles")
	m.Conf.DownloadDir = live.ParamString(p, "DownloadDir")
	if err := m.Conf.WriteConf(); err != nil {
		m.worker.AppendLog("Error saving config file. "+err.Error(), "warn")
	}
	m.ShowConf = false
	log.Printf("%+v, %+v", m.Conf, p)
	return m, nil
}

func onHtmlSnippetSave(ctx context.Context, s *live.Socket, p map[string]interface{}) (interface{}, error) {
	m := newModel(s)
	m.ShowHtmlSnippet = false

	pr, err := htmlextract.NewParseReq(p)
	if err != nil {
		m.worker.AppendLog("error extracting snippet links. "+err.Error(), "warn")
		return m, nil
	}
	links, err := htmlextract.ParseHtmlLinks(*pr)
	if err != nil {
		m.worker.AppendLog("error extracting snippet links. "+err.Error(), "warn")
		return m, nil
	}
	m.worker.AppendLog(fmt.Sprint("links in html snippet: ", len(links)), "info")
	if len(links) > 0 {
		linksLines := strings.Join(links, "\n")
		m.Links = linksLines
		m.LinksName = fmt.Sprintf("links%v", time.Now().Unix())

	}
	return m, nil
}
