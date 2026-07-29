package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jessevdk/go-flags"
)

var version string
var commit string

const UNKNOWN = 3
const CRITICAL = 2
const WARNING = 1
const OK = 0

type Opt struct {
	Version   bool          `short:"v" long:"version" description:"Show version"`
	Protocol  string        `long:"protocol" required:"true" default:"udp" choice:"tcp" choice:"udp"`
	Port      string        `short:"p" long:"port" default:"53" description:"Port number"`
	Hosts     []string      `short:"H" long:"hostname" default:"127.0.0.1" description:"DNS server hostnames"`
	Question  string        `short:"Q" long:"question" default:"example.com." description:"Question hostname"`
	QueryType string        `short:"q" long:"querytype" required:"true" default:"A" choice:"A" choice:"AAAA"`
	Expect    string        `short:"E" long:"expect" default:"" description:"Expect string in result"`
	Timeout   time.Duration `long:"timeout" default:"5s" description:"Timeout"`
	All       bool          `long:"all" description:"Require all resolution OK"`
}

func main() {
	opt := &Opt{}
	psr := flags.NewParser(opt, flags.HelpFlag|flags.PrintErrors|flags.PassDoubleDash)
	_, err := psr.Parse()
	if flags.WroteHelp(err) {
		os.Exit(OK)
	} else if err != nil {
		os.Exit(UNKNOWN)
	}
	if opt.Version {
		if commit == "" {
			commit = "dev"
		}
		fmt.Printf(
			"%s-%s\n%s/%s, %s, %s\n",
			filepath.Base(os.Args[0]),
			version,
			runtime.GOOS,
			runtime.GOARCH,
			runtime.Version(),
			commit)
		os.Exit(OK)
	}

	ckr := opt.Resolve()
	ckr.Name = "DNS-Multi"
	ckr.Exit()
}
