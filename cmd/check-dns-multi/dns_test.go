package main

import (
	"testing"
	"time"

	"github.com/mackerelio/checkers"
	"github.com/monitoring-forge/flagrun"
	"github.com/stretchr/testify/assert"
)

func TestResolveOne(t *testing.T) {
	timeout, _ := time.ParseDuration("5s")
	opt := &Opt{
		Protocol:  "udp",
		Port:      "53",
		Hosts:     []string{"8.8.8.8"},
		Question:  "example.com.",
		QueryType: "A",
		Timeout:   timeout,
	}
	msg, err := opt.resolveOne("8.8.8.8")
	if err != nil {
		t.Error(err)
	}
	assert.Contains(t, msg, "[8.8.8.8:53] ANSWER")
}

func TestResolve(t *testing.T) {
	timeout, _ := time.ParseDuration("5s")
	opt := &Opt{
		Protocol:  "udp",
		Port:      "53",
		Hosts:     []string{"8.8.8.8", "dnstestdnstestdnstest"},
		Question:  "example.com.",
		QueryType: "A",
		Timeout:   timeout,
		All:       false,
	}
	checker := opt.Run(nil)
	assert.Equal(t, checker.Status, checkers.OK)
	assert.Contains(t, checker.Message, "[8.8.8.8:53] ANSWER")
	assert.Contains(t, checker.Message, "[dnstestdnstestdnstest:53] failed to resolve")
}

func TestResolveAny(t *testing.T) {
	timeout, _ := time.ParseDuration("5s")
	opt := &Opt{
		Protocol:  "udp",
		Port:      "53",
		Hosts:     []string{"dnstestdnstestdnstest", "dnstest2dnstest2dnstest2"},
		Question:  "example.com.",
		QueryType: "A",
		Timeout:   timeout,
		All:       true,
	}
	checker := opt.Run(nil)
	assert.Equal(t, checker.Status, checkers.CRITICAL)
	assert.Contains(t, checker.Message, "[dnstestdnstestdnstest:53] failed to resolve")
	assert.Contains(t, checker.Message, "[dnstest2dnstest2dnstest2:53] failed to resolve")
}

func TestResolveAll(t *testing.T) {
	timeout, _ := time.ParseDuration("5s")
	opt := &Opt{
		Protocol:  "udp",
		Port:      "53",
		Hosts:     []string{"8.8.8.8", "dnstestdnstestdnstest"},
		Question:  "example.com.",
		QueryType: "A",
		Timeout:   timeout,
		All:       true,
	}
	checker := opt.Run(nil)
	assert.Equal(t, checker.Status, checkers.CRITICAL)
	assert.Contains(t, checker.Message, "[8.8.8.8:53] ANSWER")
	assert.Contains(t, checker.Message, "[dnstestdnstestdnstest:53] failed to resolve")
}

func TestResolveOnceExpect(t *testing.T) {
	genOpt := func(expect string, expectRegexp flagrun.Re) *Opt {
		op := &Opt{
			Protocol:  "udp",
			Port:      "53",
			Question:  "example.com.",
			QueryType: "A",
			Timeout:   5 * time.Second,
		}
		if expect != "" {
			op.Expect = expect
		}
		if expectRegexp.Present() {
			op.ExpectRegexp = expectRegexp
		}
		return op
	}
	opt := genOpt("127.0.0.1", flagrun.Re{})
	_, err := opt.resolveOne("8.8.8.8")
	assert.Error(t, err)

	opt = genOpt("", flagrun.Re{})
	err = opt.ExpectRegexp.UnmarshalFlag(`\d+\.\d+\.\d+\.\d+`)
	assert.NoError(t, err)
	_, err = opt.resolveOne("8.8.8.8")
	assert.NoError(t, err)

	opt = genOpt("", flagrun.Re{})
	err = opt.ExpectRegexp.UnmarshalFlag(`\d+\.\d+\.\d+\.A`)
	assert.NoError(t, err)
	_, err = opt.resolveOne("8.8.8.8")
	assert.Error(t, err)
}

func TestExpectAny(t *testing.T) {
	assert.True(t, expectAny("hello world", []string{"hello world"}))
	assert.True(t, expectAny("hello", []string{"hello world"}))
	assert.True(t, expectAny("hello world", []string{"world", "hello world"}))
	assert.True(t, expectAny("world", []string{"world", "hello world"}))
	assert.False(t, expectAny("hello world", []string{"foo"}))
	assert.False(t, expectAny("hello world", []string{}))
	assert.False(t, expectAny("", []string{"hello"}))
}

func TestExpectMatchAny(t *testing.T) {
	reg := flagrun.Re{}
	err := reg.UnmarshalFlag("hello")
	assert.NoError(t, err)
	assert.True(t, expectMatchAny(reg, []string{"hello world"}))
	assert.True(t, expectMatchAny(reg, []string{"world", "hello world"}))
	assert.False(t, expectMatchAny(reg, []string{"world"}))
	assert.False(t, expectMatchAny(reg, []string{}))
	assert.False(t, expectMatchAny(flagrun.Re{}, []string{"hello"}))
}
