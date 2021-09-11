package worker

import (
	"errors"
	"log"
	"os"
	"path"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"github.com/raviraa/htmltoebook/config"
)

func (w *Worker) FetchUrlChrome(url string) (string, error) {
	var e proto.NetworkResponseReceived
	var resCode, htm string

	{
		err := rod.Try(func() {
			page := w.browser.MustPage("")
			// page = page.Timeout(time.Second*idletimeout*2)
			wait := page.WaitEvent(&e)
			page.MustNavigate(url)
			wait()
			page.Mouse.MustScroll(0, 10)
			page.WaitLoad()
			// page = page.MustWaitIdle()
			page.WaitIdle(config.IdleTimeout)
			// waitidle := page.WaitRequestIdle(time.Second*idletimeout, nil, nil)
			// waitidle()

			page.WaitLoad()
			resCode = utils.Dump(e.Response.Status)
			os.WriteFile(
				path.Join(os.TempDir(), "rodhtmlpost"),
				[]byte(page.MustHTML()), 0755)
			htm = page.MustHTML()
			page.MustClose()
		})
		if err != nil {
			return "", err
		}
	}
	// fmt.Println(url, len(htm), resCode, resCode == "200")
	if resCode != "200" {
		return "", errors.New("Http Response error: " + resCode)
	}
	return htm, nil
}

func (w *Worker) stopChrome() {
	if w.browser != nil {
		w.browser.Close()
		w.launcher.Kill()
		os.RemoveAll(launcher.DefaultUserDataDirPrefix)
	}
}

func (w *Worker) resetChrome() {
	if w.conf.ChromeDownload {
		w.stopChrome()
		w.newBrowser()
	}
}

func (w *Worker) newBrowser() {
	// u := launcher.MustResolveURL("")
	// browser := rod.New().ControlURL(u).MustConnect()
	// return browser

	log.Println("Starting chrome")
	l := launcher.New().NoSandbox(true).Headless(true).Leakless(false)
	w.launcher = l

	w.browser = rod.New().Trace(false).ControlURL(l.MustLaunch()).MustConnect()

	// return rod.New().MustConnect()
	// w.browser = rod.New().MustConnect()
}
