package main

import (
	"log"
	"os"

	"github.com/afq984/ilmsmock/internal/cmdmain"
	"github.com/afq984/ilmsmock/internal/proxy"
	"github.com/spf13/pflag"
)

var upstream string
var captureRoot string = "cache"

func init() {
	pflag.StringVar(&upstream, "upstream", upstream, "")
	pflag.StringVar(&captureRoot, "capture-dir", captureRoot, "")
}

func main() {
	pflag.Parse()

	var up proxy.Upstream
	if upstream == "" {
		log.Println("using local upstream:", captureRoot)
		up = proxy.NewFileSystemUpstream(os.DirFS(captureRoot))
	} else {
		log.Println("using remote upstream:", upstream)
		up = proxy.NewHTTPUpstream(upstream, captureRoot)
	}

	cmdmain.Main(up)
}
