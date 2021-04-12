package worker

import (
	"fmt"
	"localhost/htmltoebook/config"
	"localhost/htmltoebook/types"
	"log"
	"os"
	"strings"

	"github.com/jfyne/live"
)

// Worker handles fetching html links, stripping content and conversion to mobi format
type Worker struct {
	handler *live.Handler
	done    chan bool
	running bool
}

var worker *Worker = &Worker{}

func SetHandler(h *live.Handler) {
	worker.handler = h
}

func StartWorker(links []string) {
	os.Mkdir(config.Tmpdir, 0750)
	go func() {
		log.Println("in start worker")
		out(fmt.Sprintf("Processing links: %v ", len(links)))
		FetchStripUrls(links)
		WriteMobi()
		// TODO remove title and html files.
		// appendLog(worker.handler, fmt.Sprintf("log %v %v", i, i*2), "info")
	}()

}

func AppendLog(msg, level string) {
	logmsg := types.LogMsg{Msg: msg, Level: level}
	worker.handler.Broadcast(live.Event{
		T: types.Evlogmsg,
		Data: map[string]interface{}{
			"logmsg": logmsg,
		},
	})
}

func ClearTmpDir() {
	if err := os.RemoveAll(config.Tmpdir); err != nil {
		out("Error removing temporary directory: ", config.Tmpdir, err.Error())
		return
	}
	out("Removed temporary directory: ", config.Tmpdir)
}

// log to websocket in future?
func out(s ...string) {
	log.Println(s)
	// TODO check webui running?
	AppendLog(strings.Join(s, " "), "info")
}
func panicerr(err error) {
	if err != nil {
		panic(err)
	}
}
