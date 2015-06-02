package raymond

import "testing"

//
// Those tests come from:
//   https://github.com/wycats/handlebars.js/blob/master/spec/subexpression.js
//
var hbSubexpressionsTests = []raymondTest{
	{
		"arg-less helper",
		"{{foo (bar)}}!",
		map[string]interface{}{},
		nil,
		map[string]Helper{
			"foo": func(h *HelperArg) interface{} {
				return h.ParamStr(0) + h.ParamStr(0)
			},
			"bar": func(h *HelperArg) interface{} {
				return "LOL"
			},
		},
		nil,
		"LOLLOL!",
	},
	{
		"helper w args",
		"{{blog (equal a b)}}",
		map[string]interface{}{"bar": "LOL"},
		nil,
		map[string]Helper{
			"blog":  blogHelper,
			"equal": equalHelper,
		},
		nil,
		"val is true",
	},
	{
		"mixed paths and helpers",
		"{{blog baz.bat (equal a b) baz.bar}}",
		map[string]interface{}{"bar": "LOL", "baz": map[string]string{"bat": "foo!", "bar": "bar!"}},
		nil,
		map[string]Helper{
			"blog": func(h *HelperArg) interface{} {
				return "val is " + h.ParamStr(0) + ", " + h.ParamStr(1) + " and " + h.ParamStr(2)
			},
			"equal": equalHelper,
		},
		nil,
		"val is foo!, true and bar!",
	},
	{
		"supports much nesting",
		"{{blog (equal (equal true true) true)}}",
		map[string]interface{}{"bar": "LOL"},
		nil,
		map[string]Helper{
			"blog":  blogHelper,
			"equal": equalHelper,
		},
		nil,
		"val is true",
	},

	{
		"GH-800 : Complex subexpressions (1)",
		"{{dash 'abc' (concat a b)}}",
		map[string]interface{}{"a": "a", "b": "b", "c": map[string]string{"c": "c"}, "d": "d", "e": map[string]string{"e": "e"}},
		nil,
		map[string]Helper{"dash": dashHelper, "concat": concatHelper},
		nil,
		"abc-ab",
	},
	{
		"GH-800 : Complex subexpressions (2)",
		"{{dash d (concat a b)}}",
		map[string]interface{}{"a": "a", "b": "b", "c": map[string]string{"c": "c"}, "d": "d", "e": map[string]string{"e": "e"}},
		nil,
		map[string]Helper{"dash": dashHelper, "concat": concatHelper},
		nil,
		"d-ab",
	},
	{
		"GH-800 : Complex subexpressions (3)",
		"{{dash c.c (concat a b)}}",
		map[string]interface{}{"a": "a", "b": "b", "c": map[string]string{"c": "c"}, "d": "d", "e": map[string]string{"e": "e"}},
		nil,
		map[string]Helper{"dash": dashHelper, "concat": concatHelper},
		nil,
		"c-ab",
	},
	{
		"GH-800 : Complex subexpressions (4)",
		"{{dash (concat a b) c.c}}",
		map[string]interface{}{"a": "a", "b": "b", "c": map[string]string{"c": "c"}, "d": "d", "e": map[string]string{"e": "e"}},
		nil,
		map[string]Helper{"dash": dashHelper, "concat": concatHelper},
		nil,
		"ab-c",
	},
	{
		"GH-800 : Complex subexpressions (5)",
		"{{dash (concat a e.e) c.c}}",
		map[string]interface{}{"a": "a", "b": "b", "c": map[string]string{"c": "c"}, "d": "d", "e": map[string]string{"e": "e"}},
		nil,
		map[string]Helper{"dash": dashHelper, "concat": concatHelper},
		nil,
		"ae-c",
	},

	{
		// note: test not relevant
		"provides each nested helper invocation its own options hash",
		"{{equal (equal true true) true}}",
		map[string]interface{}{},
		nil,
		map[string]Helper{
			"equal": equalHelper,
		},
		nil,
		"true",
	},
	{
		"with hashes",
		"{{blog (equal (equal true true) true fun='yes')}}",
		map[string]interface{}{"bar": "LOL"},
		nil,
		map[string]Helper{
			"blog":  blogHelper,
			"equal": equalHelper,
		},
		nil,
		"val is true",
	},
	{
		"as hashes",
		"{{blog fun=(equal (blog fun=1) 'val is 1')}}",
		map[string]interface{}{},
		nil,
		map[string]Helper{
			"blog": func(h *HelperArg) interface{} {
				return "val is " + h.HashStr("fun")
			},
			"equal": equalHelper,
		},
		nil,
		"val is true",
	},
	{
		"multiple subexpressions in a hash",
		`{{input aria-label=(t "Name") placeholder=(t "Example User")}}`,
		map[string]interface{}{},
		nil,
		map[string]Helper{
			"input": func(h *HelperArg) interface{} {
				return SafeString(`<input aria-label="` + h.HashStr("aria-label") + `" placeholder="` + h.HashStr("placeholder") + `" />`)
			},
			"t": func(h *HelperArg) interface{} {
				return SafeString(h.ParamStr(0))
			},
		},
		nil,
		`<input aria-label="Name" placeholder="Example User" />`,
	},
	{
		"multiple subexpressions in a hash with context",
		`{{input aria-label=(t item.field) placeholder=(t item.placeholder)}}`,
		map[string]map[string]string{"item": {"field": "Name", "placeholder": "Example User"}},
		nil,
		map[string]Helper{
			"input": func(h *HelperArg) interface{} {
				return SafeString(`<input aria-label="` + h.HashStr("aria-label") + `" placeholder="` + h.HashStr("placeholder") + `" />`)
			},
			"t": func(h *HelperArg) interface{} {
				return SafeString(h.ParamStr(0))
			},
		},
		nil,
		`<input aria-label="Name" placeholder="Example User" />`,
	},

	// @todo "in string params mode"

	// @todo "as hashes in string params mode"

	{
		"subexpression functions on the context",
		"{{foo (bar)}}!",
		map[string]interface{}{"bar": func() string { return "LOL" }},
		nil,
		map[string]Helper{
			"foo": func(h *HelperArg) interface{} {
				return h.ParamStr(0) + h.ParamStr(0)
			},
		},
		nil,
		"LOLLOL!",
	},

	// @todo "subexpressions can't just be property lookups" should raise error
}

func TestHandlebarsSubexpressions(t *testing.T) {
	launchHandlebarsTests(t, hbSubexpressionsTests)
}
