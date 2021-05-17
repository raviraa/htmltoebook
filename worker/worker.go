package worker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/raviraa/htmltoebook/config"
	"github.com/raviraa/htmltoebook/types"

	"github.com/jfyne/live"
)

// Worker handles fetching html links, stripping content and conversion to mobi format
type Worker struct {
	handler  *live.Handler
	conf     *config.ConfigType
	client   *http.Client
	imgCount int
}

// indicates an image in titles file
const ADDIMAGE = "ADDIMAGE"

func New(h *live.Handler, c *config.ConfigType) *Worker {
	return &Worker{
		handler: h,
		conf:    c,
		client:  &http.Client{Timeout: time.Second * 90},
	}
}

func (w *Worker) StartWorker(ctx context.Context, links []string) {
	os.MkdirAll(w.conf.Tmpdir, 0750)
	go func() {
		log.Println("in start worker")
		w.loginfo(fmt.Sprintf("Processing links: %v ", len(links)))
		if w.FetchStripUrls(ctx, links) {
			w.WriteBook()
		}
		w.notifyWebUiStop()
	}()

}

func (w *Worker) notifyWebUiStop() {
	w.handler.Broadcast(live.Event{
		T: types.EvWorkerStopped,
	})
}

func (w *Worker) AppendLog(msg, level string) {
	if w.handler == nil {
		return
	}
	logmsg := types.LogMsg{Msg: msg, Level: level}
	w.handler.Broadcast(live.Event{
		T: types.Evlogmsg,
		Data: map[string]interface{}{
			"logmsg": logmsg,
		},
	})
}

func (w *Worker) ClearTmpDir() {
	fnames, _, err := parseTitlesFile(w.conf.TitlesFname())
	log.Println(err, fnames)
	if err != nil || len(fnames) == 0 {
		w.loginfo("No intermediate files to clean up")
		return
	}
	cleaned := 0
	for _, fname := range fnames {
		if os.Remove(fname) == nil {
			cleaned++
		}
	}
	os.Remove(w.conf.TitlesFname())
	w.loginfo(fmt.Sprintf("Removed %d/%d files", cleaned, len(fnames)))
}

func (w *Worker) logsuccess(s ...string) {
	w.loglog("success", s...)
}

func (w *Worker) logerr(s ...string) {
	w.loglog("warn", s...)
}

func (w *Worker) loginfo(s ...string) {
	w.loglog("info", s...)
}

// loglog logs to both the ui and cli
func (w *Worker) loglog(level string, logs ...string) {
	// find calling function(2 levels deep) file name and line number
	_, cfile, cline, _ := runtime.Caller(2)
	cfileSpl := strings.Split(cfile, "/")
	caller := fmt.Sprintf("%s:%d", cfileSpl[len(cfileSpl)-1], cline)
	timenow := time.Now().Format("15:04:05 ")
	logsjoin := strings.Join(logs, " ")

	fmt.Println(timenow, caller, logsjoin)
	if w.handler != nil {
		w.AppendLog(timenow+logsjoin, level)
	}
}
