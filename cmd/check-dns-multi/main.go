package main

import (
	"os"
	"time"

	"github.com/monitoring-forge/flagrun"
)

var version string

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
	os.Exit(flagrun.Check(&Opt{}, flagrun.Version(version)))
}
