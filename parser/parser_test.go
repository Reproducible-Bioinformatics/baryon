package parser

import "testing"

func Test_RlangDocParser(t *testing.T) {
}

func Test_IsRoxygenLine(t *testing.T) {
	type testStruct struct {
		Expect bool
		Line   string
	}
	var tests = []testStruct{
		{Expect: false, Line: "Not a roxygen line"},
		{Expect: true, Line: "#' A roxygen line"},
		{Expect: true, Line: "#'"},
		{Expect: false, Line: "#"},
	}

	for _, entry := range tests {
		if entry.Expect != isRoxygenLine(entry.Line) {
			t.Errorf("%s got a wrong expectation?", entry.Line)
		}
	}
}
