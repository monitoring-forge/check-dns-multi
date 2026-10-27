package main

import (
	"testing"

	"github.com/monitoring-forge/flagrun"
)

func TestOptValidator(t *testing.T) {
	testFlagRegexp := flagrun.Re{}
	err := testFlagRegexp.UnmarshalFlag("^example")
	if err != nil {
		t.Fatalf("failed to unmarshal flag: %v", err)
	}
	testCases := []struct {
		expect       string
		expectRegexp flagrun.Re
		err          bool
	}{
		{"", flagrun.Re{}, false},
		{"value", flagrun.Re{}, false},
		{"", testFlagRegexp, false},
		{"value", testFlagRegexp, true},
	}
	for _, tc := range testCases {
		opt := &Opt{
			Expect:       tc.expect,
			ExpectRegexp: tc.expectRegexp,
		}
		err := opt.Validator(nil)
		if !tc.err && err != nil {
			t.Errorf("unexpected error for expect %q and expect-regexp %q: %v", tc.expect, tc.expectRegexp.String(), err)
		}
		if tc.err && err == nil {
			t.Errorf("expected error for expect %q and expect-regexp %q, but got none", tc.expect, tc.expectRegexp.String())
		}
	}
}

func TestEmptyOptValiddator(t *testing.T) {
	opt := &Opt{}
	err := opt.Validator(nil)
	if err != nil {
		t.Errorf("unexpected error for empty Opt: %v", err)
	}
}
