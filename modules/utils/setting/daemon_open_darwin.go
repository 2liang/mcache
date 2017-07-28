package setting

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/facebookgo/grace/gracehttp"
)

var (
	IsDebug bool = false
)

func Daemon() {

	debug := flag.Bool("debug", false, "debug mode ?")
	flag.Parse()
	IsDebug = *debug
	if os.Getenv("LISTEN_FDS") == "" && os.Getppid() != 1 && !*debug {
		var se string = os.Args[0]
		var args []string = os.Args[1:]
		filter(args)
		cmd := exec.Command(se, args...)
		err := cmd.Start()

		if err != nil {
			Logger.Panic("Start unsuccessfully")
		}
		os.Exit(0)
	}
}

//Filter the flag of daemon
func filter(args []string) {

	for arg := range args {
		if strings.EqualFold(args[arg], "--debug") || strings.EqualFold(args[arg], "-debug") {
			args[arg] = ""
			return
		}
	}

}

//initialize the signal controal
func SetSignal(serve *http.Server) {

	/*var c chan os.Signal
	c = make(chan os.Signal)

	go func() {
		signal.Notify(c, os.Interrupt, os.Kill)
	}()

	for {
		switch <-c {
		case os.Interrupt:
			Logger.Fatalf("%s", "The GOil was interrupted")
			return
		case os.Kill:
			Logger.Fatalf("%s", "The GOil was killed")
			return
		}
	}*/
	gracehttp.Serve(serve)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	serve.Shutdown(ctx)
	Logger.Info("Server already shutdown gracefully!")
	//shutdown gracefully
	/*select {
	case <-ctx.Done():
		Logger.Info("Server already shutdown gracefully!")
	}*/
}
