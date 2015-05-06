package parser

import (
	"log"
	"testing"

	"github.com/aymerick/raymond/ast"
)

const (
	VERBOSE = false
)

type parserTest struct {
	name   string
	input  string
	output string
}

var parserTests = []parserTest{
	{"Content", "Hello", "CONTENT[Hello]\n"},
	{"Comment", "{{! This is a comment }}", "{{! 'This is a comment' }}\n"},
	{"Comment dash", "{{!-- This is a comment --}}", "{{! 'This is a comment' }}\n"},

	//
	// Next tests come from:
	//   https://github.com/wycats/handlebars.js/blob/master/spec/parser.js
	//
	{"parses simple mustaches (1)", `{{123}}`, "{{ NUMBER{123} [] }}\n"},
	{"parses simple mustaches (2)", `{{"foo"}}`, "{{ \"foo\" [] }}\n"},
	{"parses simple mustaches (3)", `{{false}}`, "{{ BOOLEAN{false} [] }}\n"},
	{"parses simple mustaches (4)", `{{true}}`, "{{ BOOLEAN{true} [] }}\n"},
	{"parses simple mustaches (5)", `{{foo}}`, "{{ PATH:foo [] }}\n"},
	{"parses simple mustaches (6)", `{{foo?}}`, "{{ PATH:foo? [] }}\n"},
	{"parses simple mustaches (7)", `{{foo_}}`, "{{ PATH:foo_ [] }}\n"},
	{"parses simple mustaches (8)", `{{foo-}}`, "{{ PATH:foo- [] }}\n"},
	{"parses simple mustaches (9)", `{{foo:}}`, "{{ PATH:foo: [] }}\n"},

	// {"parses simple mustaches with data", `{{@foo}}`, "{{ @PATH:foo [] }}\n"},
	// {"parses simple mustaches with data paths", `{{@../foo}}`, "{{ @PATH:foo [] }}\n"},
	// {"parses mustaches with paths", `{{foo/bar}}`, "{{ PATH:foo/bar [] }}\n"},
	// {"parses mustaches with this/foo", `{{this/foo}}`, "{{ PATH:foo [] }}\n"},
	// {"parses mustaches with - in a path", `{{foo-bar}}`, "{{ PATH:foo-bar [] }}\n"},
	// {"parses mustaches with parameters", `{{foo bar}}`, "{{ PATH:foo [PATH:bar] }}\n"},
	// {"parses mustaches with string parameters", `{{foo bar \"baz\" }}`, "{{ PATH:foo [PATH:bar, \"baz\"] }}\n"},
	// {"parses mustaches with NUMBER parameters", `{{foo 1}}`, "{{ PATH:foo [NUMBER{1}] }}\n"},
	// {"parses mustaches with BOOLEAN parameters (1)", `{{foo true}}`, "{{ PATH:foo [BOOLEAN{true}] }}\n"},
	// {"parses mustaches with BOOLEAN parameters (2)", `{{foo false}}`, "{{ PATH:foo [BOOLEAN{false}] }}\n"},
	// {"parses mutaches with DATA parameters", `{{foo @bar}}`, "{{ PATH:foo [@PATH:bar] }}\n"},
}

func TestParser(t *testing.T) {
	for _, test := range parserTests {
		if VERBOSE {
			log.Printf("\n\n**********************************")
			log.Printf("Testing: %s", test.name)
		}

		output := ""

		node, err := Parse(test.input)
		if err == nil {
			output = ast.PrintNode(node)
		}

		if (err != nil) || (test.output != output) {
			t.Errorf("Test '%s' failed\ninput:\n\t'%s'\nexpected\n\t%q\ngot\n\t%q\nerror:\n\t%s", test.name, test.input, test.output, output, err)
		}
	}
}
