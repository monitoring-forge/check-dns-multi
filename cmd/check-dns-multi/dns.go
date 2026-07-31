package main

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/mackerelio/checkers"
	"github.com/miekg/dns"
)

func (o *Opt) resolveOne(host string) (string, error) {
	c := &dns.Client{
		Net:     o.Protocol,
		Timeout: o.Timeout,
	}
	address := net.JoinHostPort(host, o.Port)
	m := new(dns.Msg)
	m.SetQuestion(o.Question, dns.StringToType[o.QueryType])
	r, _, err := c.Exchange(m, address)
	if err != nil {
		return "", fmt.Errorf("[%s] failed to resolve: %w", address, err)
	}
	if r.Rcode != dns.RcodeSuccess {
		return "", fmt.Errorf("[%s] failed to resolve '%s': rcode:%s",
			address,
			o.Question,
			dns.RcodeToString[r.Rcode],
		)
	}
	answer := make([]string, 0)
	for _, a := range r.Answer {
		if aa, ok := a.(*dns.A); ok {
			answer = append(answer, aa.A.String())
		}
		if aa, ok := a.(*dns.AAAA); ok {
			answer = append(answer, aa.AAAA.String())
		}
	}
	if len(o.Expect) > 0 && !strings.Contains(strings.Join(answer, "|"), o.Expect) {
		return "", fmt.Errorf("[%s] not contain '%s' in '%s'",
			address,
			o.Expect,
			strings.Join(answer, "|"),
		)
	}
	msg := make([]string, 0, len(r.Answer)+1)
	msg = append(msg, fmt.Sprintf("[%s] HEADER-> %s", address,
		strings.ReplaceAll(r.MsgHdr.String(), "\n", " ")))
	for _, answer := range r.Answer {
		msg = append(msg,
			strings.ReplaceAll(fmt.Sprintf("[%s] ANSWER-> %s", address, answer), "\n", " "))
	}
	return strings.Join(msg, "\n"), nil
}

func (o *Opt) Resolve() *checkers.Checker {
	var mu sync.Mutex
	errCount := 0
	m := make([]string, 0, len(o.Hosts))

	var wg sync.WaitGroup
	for _, host := range o.Hosts {
		wg.Go(func() {
			msg, err := o.resolveOne(host)
			mu.Lock()
			if err != nil {
				errCount++
				m = append(m, err.Error())
			} else {
				m = append(m, msg)
			}
			mu.Unlock()
		})
	}
	wg.Wait()

	checkSt := checkers.OK
	if o.All {
		// should be errCount == 0
		if errCount > 0 {
			checkSt = checkers.CRITICAL
		}
	} else {
		// err if error all
		if errCount == len(o.Hosts) {
			checkSt = checkers.CRITICAL
		}
	}
	return checkers.NewChecker(checkSt, strings.Join(m, "\n"))
}
