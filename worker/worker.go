package worker

import (
	"fmt"
	"localhost/htmltoebook/types"
	"log"
	"time"

	"github.com/jfyne/live"
)

// Worker handles fetching html links, stripping content and conversion to mobi format
type Worker struct {
	handler *live.Handler
	done    chan bool
	running bool
}

// Worker singleton
var worker *Worker

func init() {
	worker = &Worker{}
}

func SetHandler(h *live.Handler) {
	worker.handler = h
}

func StartWorker() {
	// TODO check is running?
	go func() {
		log.Println("in start worker")
		for i := 0; i < 9; i++ {
			time.Sleep(100 * time.Millisecond)
			appendLog(worker.handler, fmt.Sprintf("log %v %v", i, i*2), "info")
		}
	}()

}

func appendLog(h *live.Handler, msg, level string) {
	logmsg := types.LogMsg{Msg: msg, Level: level}
	h.Broadcast(live.Event{
		T: types.Evlogmsg,
		Data: map[string]interface{}{
			"logmsg": logmsg,
		},
	})
}
