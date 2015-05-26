package raymond

import "testing"

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/whitespace-control.js
//
var hbWhitespaceControlTests = []raymondTest{
	{
		"should strip whitespace around mustache calls (1)",
		" {{~foo~}} ",
		map[string]string{"foo": "bar<"},
		nil,
		"bar&lt;",
	},
	{
		"should strip whitespace around mustache calls (2)",
		" {{~foo}} ",
		map[string]string{"foo": "bar<"},
		nil,
		"bar&lt; ",
	},
	{
		"should strip whitespace around mustache calls (3)",
		" {{foo~}} ",
		map[string]string{"foo": "bar<"},
		nil,
		" bar&lt;",
	},
	{
		"should strip whitespace around mustache calls (4)",
		" {{~&foo~}} ",
		map[string]string{"foo": "bar<"},
		nil,
		"bar<",
	},
	{
		"should strip whitespace around mustache calls (5)",
		" {{~{foo}~}} ",
		map[string]string{"foo": "bar<"},
		nil,
		"bar<",
	},
	{
		"should strip whitespace around mustache calls (6)",
		"1\n{{foo~}} \n\n 23\n{{bar}}4",
		nil,
		nil,
		"1\n23\n4",
	},

	{
		"blocks - should strip whitespace around simple block calls (1)",
		" {{~#if foo~}} bar {{~/if~}} ",
		map[string]string{"foo": "bar<"},
		nil,
		"bar",
	},
	{
		"blocks - should strip whitespace around simple block calls (2)",
		" {{#if foo~}} bar {{/if~}} ",
		map[string]string{"foo": "bar<"},
		nil,
		" bar ",
	},
	{
		"blocks - should strip whitespace around simple block calls (3)",
		" {{~#if foo}} bar {{~/if}} ",
		map[string]string{"foo": "bar<"},
		nil,
		" bar ",
	},
	{
		"blocks - should strip whitespace around simple block calls (4)",
		" {{#if foo}} bar {{/if}} ",
		map[string]string{"foo": "bar<"},
		nil,
		"  bar  ",
	},
	{
		"blocks - should strip whitespace around simple block calls (5)",
		" \n\n{{~#if foo~}} \n\nbar \n\n{{~/if~}}\n\n ",
		map[string]string{"foo": "bar<"},
		nil,
		"bar",
	},
	{
		"blocks - should strip whitespace around simple block calls (6)",
		" a\n\n{{~#if foo~}} \n\nbar \n\n{{~/if~}}\n\na ",
		map[string]string{"foo": "bar<"},
		nil,
		" abara ",
	},

	{
		"should strip whitespace around inverse block calls (1)",
		" {{~^if foo~}} bar {{~/if~}} ",
		nil,
		nil,
		"bar",
	},
	{
		"should strip whitespace around inverse block calls (2)",
		" {{^if foo~}} bar {{/if~}} ",
		nil,
		nil,
		" bar ",
	},
	{
		"should strip whitespace around inverse block calls (3)",
		" {{~^if foo}} bar {{~/if}} ",
		nil,
		nil,
		" bar ",
	},
	{
		"should strip whitespace around inverse block calls (4)",
		" {{^if foo}} bar {{/if}} ",
		nil,
		nil,
		"  bar  ",
	},
	{
		"should strip whitespace around inverse block calls (5)",
		" \n\n{{~^if foo~}} \n\nbar \n\n{{~/if~}}\n\n ",
		nil,
		nil,
		"bar",
	},

	// {
	// 	"should strip whitespace around complex block calls (1)",
	// 	"{{#if foo~}} bar {{~^~}} baz {{~/if}}",
	// 	map[string]string{"foo": "bar<"},
	// 	nil,
	// 	"bar",
	// },
	// {
	// 	"should strip whitespace around complex block calls (2)",
	// 	"{{#if foo~}} bar {{^~}} baz {{/if}}",
	// 	map[string]string{"foo": "bar<"},
	// 	nil,
	// 	"bar ",
	// },
	// {
	// 	"should strip whitespace around complex block calls (3)",
	// 	"{{#if foo}} bar {{~^~}} baz {{~/if}}",
	// 	map[string]string{"foo": "bar<"},
	// 	nil,
	// 	" bar",
	// },
	// {
	// 	"should strip whitespace around complex block calls (4)",
	// 	"{{#if foo}} bar {{^~}} baz {{/if}}",
	// 	map[string]string{"foo": "bar<"},
	// 	nil,
	// 	" bar ",
	// },
	// {
	// 	"should strip whitespace around complex block calls (5)",
	// 	"{{#if foo~}} bar {{~else~}} baz {{~/if}}",
	// 	map[string]string{"foo": "bar<"},
	// 	nil,
	// 	"bar",
	// },
	// {
	// 	"should strip whitespace around complex block calls (6)",
	// 	"\n\n{{~#if foo~}} \n\nbar \n\n{{~^~}} \n\nbaz \n\n{{~/if~}}\n\n",
	// 	map[string]string{"foo": "bar<"},
	// 	nil,
	// 	"bar",
	// },
	// {
	// 	"should strip whitespace around complex block calls (7)",
	// 	"\n\n{{~#if foo~}} \n\n{{{foo}}} \n\n{{~^~}} \n\nbaz \n\n{{~/if~}}\n\n",
	// 	map[string]string{"foo": "bar<"},
	// 	nil,
	// 	"bar<",
	// },
	// {
	// 	"should strip whitespace around complex block calls (8)",
	// 	"{{#if foo~}} bar {{~^~}} baz {{~/if}}",
	// 	nil,
	// 	nil,
	// 	"baz",
	// },
	// {
	// 	"should strip whitespace around complex block calls (9)",
	// 	"{{#if foo}} bar {{~^~}} baz {{/if}}",
	// 	nil,
	// 	nil,
	// 	"baz ",
	// },
	// {
	// 	"should strip whitespace around complex block calls (10)",
	// 	"{{#if foo~}} bar {{~^}} baz {{~/if}}",
	// 	nil,
	// 	nil,
	// 	" baz",
	// },
	// {
	// 	"should strip whitespace around complex block calls (11)",
	// 	"{{#if foo~}} bar {{~^}} baz {{/if}}",
	// 	nil,
	// 	nil,
	// 	" baz ",
	// },
	// {
	// 	"should strip whitespace around complex block calls (12)",
	// 	"{{#if foo~}} bar {{~else~}} baz {{~/if}}",
	// 	nil,
	// 	nil,
	// 	"baz",
	// },
	// {
	// 	"should strip whitespace around complex block calls (13)",
	// 	"\n\n{{~#if foo~}} \n\nbar \n\n{{~^~}} \n\nbaz \n\n{{~/if~}}\n\n",
	// 	nil,
	// 	nil,
	// 	"baz",
	// },

	// @todo Add remaining tests
}

func TestHandlebarsWhitespaceControl(t *testing.T) {
	launchHandlebarsTests(t, hbWhitespaceControlTests)
}
