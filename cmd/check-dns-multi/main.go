package main

import (
	"fmt"
	"os"
	"time"

	"github.com/monitoring-forge/flagrun"
)

var version string

type Opt struct {
	Version      bool          `short:"v" long:"version" description:"Show version"`
	Protocol     string        `long:"protocol" required:"true" default:"udp" choice:"tcp" choice:"udp"` //nolint:staticcheck
	Port         string        `short:"p" long:"port" default:"53" description:"Port number"`
	Hosts        []string      `short:"H" long:"hostname" default:"127.0.0.1" description:"DNS server hostnames"`
	Question     string        `short:"Q" long:"question" default:"example.com." description:"Question hostname"`
	QueryType    string        `short:"q" long:"querytype" required:"true" default:"A" choice:"A" choice:"AAAA"` //nolint:staticcheck
	Expect       string        `short:"E" long:"expect" default:"" description:"Expect string in result"`
	ExpectRegexp flagrun.Re    `short:"R" long:"expect-regexp" default:"" description:"Expect regular expression in result"`
	Timeout      time.Duration `long:"timeout" default:"5s" description:"Timeout"`
	All          bool          `long:"all" description:"Require all resolution OK"`
}

func (o *Opt) Validator(_ []string) error {
	if o.ExpectRegexp.Present() && o.Expect != "" {
		return fmt.Errorf("expect-regexp is set but expect is also set")
	}
	return nil
}

func main() {
	opt := &Opt{}
	os.Exit(flagrun.Check(opt, flagrun.Version(version), flagrun.Validator(opt.Validator)))
}
