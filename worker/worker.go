package worker

import (
	"context"
	"fmt"
	"localhost/htmltoebook/config"
	"localhost/htmltoebook/types"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jfyne/live"
)

// Worker handles fetching html links, stripping content and conversion to mobi format
type Worker struct {
	handler *live.Handler
}

var worker *Worker = &Worker{}

func SetHandler(h *live.Handler) {
	worker.handler = h
}

func StartWorker(ctx context.Context, links []string) {
	os.Mkdir(config.Tmpdir, 0750)
	go func() {
		log.Println("in start worker")
		loginfo(fmt.Sprintf("Processing links: %v ", len(links)))
		if FetchStripUrls(ctx, links) {
			WriteMobi()
		}
		notifyWebUiStop()
	}()

}

func notifyWebUiStop() {
	worker.handler.Broadcast(live.Event{
		T: types.EvWorkerStopped,
	})
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
		logerr("Error removing temporary directory: ", config.Tmpdir, err.Error())
		return
	}
	loginfo("Removed temporary directory: ", config.Tmpdir)
}

func logsuccess(s ...string) {
	loglog("success", s...)
}

func logerr(s ...string) {
	loglog("warn", s...)
}

func loginfo(s ...string) {
	loglog("info", s...)
}

func loglog(level string, s ...string) {
	log.Println(s)
	if worker.handler != nil {
		outstr := time.Now().Format("15:04:05 ")
		outstr += strings.Join(s, " ")
		AppendLog(outstr, level)
	}
}
func panicerr(err error) {
	if err != nil {
		panic(err)
	}
}
