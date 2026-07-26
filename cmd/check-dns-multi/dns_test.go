package main

import (
	"testing"
	"time"

	"github.com/mackerelio/checkers"
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
	checker := opt.Resolve()
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
	checker := opt.Resolve()
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
	checker := opt.Resolve()
	assert.Equal(t, checker.Status, checkers.CRITICAL)
	assert.Contains(t, checker.Message, "[8.8.8.8:53] ANSWER")
	assert.Contains(t, checker.Message, "[dnstestdnstestdnstest:53] failed to resolve")
}
