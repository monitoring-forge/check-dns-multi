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

const (
	OK = iota
	WARNING
	CRITICAL
	UNKNOWN
)

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
	opt := Opt{}
	psr := flags.NewParser(&opt, flags.HelpFlag|flags.PassDoubleDash)
	_, err := psr.Parse()
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
	} else if flags.WroteHelp(err) {
		fmt.Fprintf(os.Stdout, "%v\n", err)
		os.Exit(OK)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(UNKNOWN)
	}

	ckr := opt.Resolve()
	ckr.Name = "DNS-Multi"
	ckr.Exit()
}
