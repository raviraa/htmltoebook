package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/raviraa/htmltoebook/config"
	"github.com/raviraa/htmltoebook/htmlextract"
	"github.com/raviraa/htmltoebook/worker"
)

const tmpfname = "/tmp/htmltoebook.toml"

const linksmsg = ` #Add web links below
 #Each link in a separate line
 #Lines that do not start with http or https will be ignored
 #Numeric range can be used: http://example.com/page{1-10}}"

https://blog.golang.org/go1.15
https://blog.golang.org/go1.16
`

func RunLinks() {
	if err := writeEditorFile(*config.New(), linksmsg); err != nil {
		log.Fatal("failed generating links template. ", err)
	}
	for {
		conf := &config.ConfigType{}
		linkslines, err := runEditor(conf)
		if err != nil {
			fmt.Println("unable to parse configuration from editor, ", err)
			if cmdAsk("Do you want to try again?") {
				continue
			}
			return
		}
		links := worker.SplitLinks(linkslines)
		conf.WriteConf()
		startWorker(conf, links)
		break
	}
}

const snippetmsg = `<p>Add html snippet code here.</p>
<p>Only anchor tags matching filters are considered</p>`

func RunHtmlSnippet() {
	if err := writeEditorFile(htmlextract.ParseReq{}, snippetmsg); err != nil {
		log.Fatal("failed generating editor template file, ", err)
	}
	for {
		prnew := &htmlextract.ParseReq{}
		htm, err := runEditor(prnew)
		if err != nil {
			fmt.Println("unable to parse configuration from editor, ", err)
			if cmdAsk("Do you want to try again?") {
				continue
			}
			return
		}

		prnew.HtmlSnippet = htm
		links, err := htmlextract.ParseHtmlLinks(*prnew)
		if err != nil {
			fmt.Println("error extracting links. ", err)
			if cmdAsk("Do you want to try again?") {
				continue
			}
			return
		}
		// log.Printf("%#v\n", *prnew)
		fmt.Printf("%s\nFound %d links\n", strings.Join(links, "\n"), len(links))
		if !cmdAsk("Proceed with the links(y), or Retry(n)") {
			continue
		}
		conf := config.New()
		startWorker(conf, links)
		break
	}
}

func RunCmd(cmdname string, args ...string) error {
	cmd := exec.Command(cmdname, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getEditor() string {
	return os.Getenv("EDITOR")
}

func cmdAsk(q string) bool {
	fmt.Print(q + " (y/n)")
	var inp string
	fmt.Fscanf(os.Stdin, "%s", &inp)
	return inp == "y"
}

func startWorker(c *config.ConfigType, links []string) {
	log.Printf("Starting with %v links\n", len(links))
	ch := make(chan bool)
	w := worker.New(nil, ch, c)
	w.StartWorker(context.Background(), links)
	<-ch
}
