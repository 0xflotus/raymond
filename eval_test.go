package raymond

import "testing"

var evalTests = []raymondTest{
	{
		"only content",
		"this is content",
		nil,
		nil,
		"this is content",
	},
	// @todo Test with a struct for data
}

func TestEval(t *testing.T) {
	launchRaymondTests(t, evalTests)
}

// @todo Test with a "../../path" (depth 2 path) while context is only depth 1

//
// StrValue() / StrInterface() tests
//

type strTest struct {
	name   string
	input  interface{}
	output string
}

var strTests = []strTest{
	{"String", "foo", "foo"},
	{"Boolean true", true, "true"},
	{"Boolean false", false, "false"},
	{"Integer", 25, "25"},
	{"Float", 25.75, "25.75"},
	{"Nil", nil, ""},
	{"[]string", []string{"foo", "bar"}, "foobar"},
	{"[]interface{} (strings)", []interface{}{"foo", "bar"}, "foobar"},
	{"[]Boolean", []bool{true, false}, "truefalse"},
}

func TestStrInterface(t *testing.T) {
	for _, test := range strTests {
		if res := StrInterface(test.input); res != test.output {
			t.Errorf("Failed to stringify: %s\nexpected:\n\t'%s'got:\n\t%q", test.name, test.output, res)
		}
	}
}
