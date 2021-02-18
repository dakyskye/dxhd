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
			data: parser{
				lineNumber: 1,
				line:       "#!/bin/sh",
				isPrefix:   false,
			},
			res: true,
		},
		{
			data: parser{
				lineNumber: 1,
				line:       "#!/bin/sh",
				isPrefix:   true,
			},
			res: false,
		},
		{
			data: parser{
				lineNumber: 0,
				line:       "#!/bin/sh",
				isPrefix:   false,
			},
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
			data: parser{line: ""},
			res:  true,
		},
		{
			data: parser{line: "  "},
			res:  true,
		},
		{
			data: parser{line: " . "},
			res:  false,
		},
	}

	for i, c := range testCases {
		if c.data.isEmptyLine() != c.res {
			t.Error(fmt.Errorf("expected result for test case %d was %t", i, c.res))
		}
	}
}
