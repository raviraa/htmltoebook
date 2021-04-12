package web

import (
	"context"
	"fmt"
	"localhost/htmltoebook/types"
	"localhost/htmltoebook/worker"
	"log"
	"net/http"

	"github.com/jfyne/live"
)

func setEvents(h *live.Handler) {
	h.Mount = onMount
	h.HandleEvent(evstart, onStart)
	h.HandleSelf(evlogmsg, onLogMsg)
	h.HandleEvent(evstop, onStop)
	h.HandleEvent(evclear, onClear)
}

func onMount(ctx context.Context, r *http.Request, s *live.Socket) (interface{}, error) {
	m := newModel(s)
	if s.Connected() {
		go func() {
			log.Println("Connected")
			worker.AppendLog("Connected. Ready to start.", "info")

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
	m.LogMsgs = []types.LogMsg{logmsg}
	return m, nil
}

func onClear(ctx context.Context, s *live.Socket, p map[string]interface{}) (interface{}, error) {
	worker.ClearTmpDir()
	return newModel(s), nil
}

func onStop(ctx context.Context, s *live.Socket, p map[string]interface{}) (interface{}, error) {
	log.Println(p)
	m := newModel(s)

	return m, nil
}

func onStart(ctx context.Context, s *live.Socket, p map[string]interface{}) (interface{}, error) {
	log.Println(p)
	// TODO check is running?
	m := newModel(s)
	linksParam := live.ParamString(p, "links")
	links := worker.SplitLinks(linksParam)
	worker.StartWorker(links)

	return m, nil
}
