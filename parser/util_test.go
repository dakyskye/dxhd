package parser

import (
	"fmt"
	"testing"
)

func TestIsShebang(t *testing.T) {
	testCases := []struct {
		data parser
		res  bool
	}{
		{
			data: parser{state: state{
				lineNumber: 1,
				line:       []byte("#!/bin/sh"),
				isPrefix:   false,
			}},
			res: true,
		},
		{
			data: parser{state: state{
				lineNumber: 1,
				line:       []byte("#!/bin/sh"),
				isPrefix:   true,
			}},
			res: false,
		},
		{
			data: parser{state: state{
				lineNumber: 0,
				line:       []byte("#!/bin/sh"),
				isPrefix:   false,
			}},
			res: false,
		},
	}

	for i, c := range testCases {
		if c.data.isShebang() != c.res {
			t.Error(fmt.Errorf("expected result for test case %d was %t", i, c.res))
		}
	}
}

func TestIsEmptyLine(t *testing.T) {
	testCases := []struct {
		data parser
		res  bool
	}{
		{
			data: parser{state: state{line: []byte("")}},
			res:  true,
		},
		{
			data: parser{state: state{line: []byte("  ")}},
			res:  true,
		},
		{
			data: parser{state: state{line: []byte(" . ")}},
			res:  false,
		},
	}

	for i, c := range testCases {
		if c.data.isEmptyLine() != c.res {
			t.Error(fmt.Errorf("expected result for test case %d was %t", i, c.res))
		}
	}
}

func TestIsAKeybinding(t *testing.T) {
	testCases := []struct {
		data parser
		res  bool
	}{
		{
			data: parser{state: state{line: []byte("#super  + foo")}},
			res:  true,
		},
		{
			data: parser{state: state{line: []byte("#super  + {foo, bar}")}},
			res:  true,
		},
		{
			data: parser{state: state{line: []byte("#super  + {foo, 1-10}")}},
			res:  true,
		},
		{
			data: parser{state: state{line: []byte("#super  + {_, baz + } {a-z, 1-10, foo, bar}")}},
			res:  true,
		},
		{
			data: parser{state: state{line: []byte("#super  + {foo bar} + {baz, qux}")}},
			res:  false,
		},
		{
			data: parser{state: state{line: []byte("super + a")}},
			res:  false,
		},
	}

	for i, c := range testCases {
		if c.data.isAKeybinding() != c.res {
			t.Error(fmt.Errorf("expected result for test case %d was %t", i, c.res))
		}
	}
}
