package parser

import (
	"reflect"
	"testing"
)

func Test_RlangDocParser(t *testing.T) {
}

func Test_IsRoxygenLine(t *testing.T) {
	type testStruct struct {
		Expect     bool
		Line       string
		Submatches []string
	}
	var tests = []testStruct{
		{Expect: false, Line: "Not a roxygen line", Submatches: []string{}},
		{Expect: true, Line: "#' A roxygen line", Submatches: []string{"A roxygen line"}},
		{Expect: true, Line: "#'", Submatches: []string{}},
		{Expect: false, Line: "#", Submatches: []string{}},
	}

	for _, entry := range tests {
		ok, subs := isRoxygenLine(entry.Line)
		if ok != entry.Expect && !reflect.DeepEqual(entry.Submatches, subs) {
			t.Errorf("%s got a wrong expectation?", entry.Line)
		}
	}
}
