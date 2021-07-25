package cmdmain

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"

	"github.com/afq984/ilmsmock/internal/proxy"
	"github.com/gorilla/handlers"
	"github.com/spf13/pflag"
)

var setEnv string
var verbose bool

func init() {
	pflag.StringVar(&setEnv, "setenv", setEnv, "")
	pflag.BoolVar(&verbose, "verbose", verbose, "turn on HTTP logging")
}

func Main(upstream proxy.Upstream) {
	if setEnv == "" {
		log.Fatal("--setenv not set")
	}

	lis, err := net.Listen("tcp4", "localhost:")
	if err != nil {
		log.Fatal(err)
	}

	var s http.Handler = &proxy.Server{
		Upstream: upstream,
	}
	if verbose {
		s = handlers.LoggingHandler(os.Stderr, s)
	}

	go func() {
		log.Fatal(http.Serve(lis, s))
	}()

	cmd := exec.Command(pflag.Arg(0), pflag.Args()[1:]...)
	cmd.Env = append(os.Environ(), setEnv+"=http://"+lis.Addr().String())
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
