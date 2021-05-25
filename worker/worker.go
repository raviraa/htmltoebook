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
	cmddone  chan bool
	imgCount int
}

// indicates an image in titles file
const ADDIMAGE = "ADDIMAGE"

func New(h *live.Handler, cmdnotif chan bool, c *config.ConfigType) *Worker {
	return &Worker{
		handler: h,
		conf:    c,
		cmddone: cmdnotif,
		client:  &http.Client{Timeout: time.Second * 90},
	}
}

func (w *Worker) StartWorker(ctx context.Context, links []string) {
	os.MkdirAll(w.conf.Tmpdir(), 0750)
	go func() {
		log.Println("Starting worker")
		w.loginfo(fmt.Sprintf("Processing links: %v ", len(links)))
		if w.FetchStripUrls(ctx, links) {
			w.WriteBook()
		}
		w.notifyWebUiStop()
	}()

}

func (w *Worker) notifyWebUiStop() {
	if w.handler != nil {
		w.handler.Broadcast(live.Event{
			T: types.EvWorkerStopped,
		})

	} else {
		w.cmddone <- true
	}
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
	os.Remove(w.conf.Tmpdir())
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

var (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	cmdColors  = map[string]string{
		"success": colorGreen,
		"warn":    colorRed,
	}
)

// loglog logs to both the ui and cli
func (w *Worker) loglog(level string, logs ...string) {
	// find calling function(2 levels deep) file name and line number
	_, cfile, cline, _ := runtime.Caller(2)
	cfileSpl := strings.Split(cfile, "/")
	caller := fmt.Sprintf("%s:%d", cfileSpl[len(cfileSpl)-1], cline)
	timenow := time.Now().Format("15:04:05")
	logsjoin := strings.Join(logs, " ")

	if color, ok := cmdColors[level]; ok {
		fmt.Println(color+timenow, caller, logsjoin, colorReset)
	} else {
		fmt.Println(timenow, caller, logsjoin)
	}
	if w.handler != nil {
		w.AppendLog(timenow+logsjoin, level)
	}
}
