package jira

import (
	//	"fmt"
	"reflect"
	"testing"
)

type test struct {
	name string
	in   []string
	out  []string
}

type testGroup map[string][]test

var singleFeatureTests = testGroup{
	"unordered list": []test{
		{
			name: "with whitespace",
			in: []string{
				"		*** 3 indent unordered",
			},
			out: []string{
				"      * 3 indent unordered",
			},
		},
		{
			name: "without whitespace",
			in: []string{
				"*** 3 indent unordered",
			},
			out: []string{
				"      * 3 indent unordered",
			},
		},
		{
			name: "nested indents",
			in: []string{
				"* level one",
				"* level one",
				"** level two",
				"** level two",
				"* level one",
			},
			out: []string{
				"  * level one",
				"  * level one",
				"    * level two",
				"    * level two",
				"  * level one",
			},
		},
	},
	"ordered list": []test{
		{
			name: "with whitespace",
			in: []string{
				"		### 3 indent ordered",
			},
			out: []string{
				"      1. 3 indent ordered",
			},
		},
		{
			name: "without whitespace",
			in: []string{
				"### 3 indent ordered",
			},
			out: []string{
				"      1. 3 indent ordered",
			},
		},
		{
			name: "nested indents",
			in: []string{
				"# level one",
				"# level one",
				"## level two",
				"## level two",
				"# level one",
			},
			out: []string{
				"  1. level one",
				"  1. level one",
				"    1. level two",
				"    1. level two",
				"  1. level one",
			},
		},
	},
	"header": []test{
		{
			name: "requires dot",
			in: []string{
				"h1 not a header",
			},
			out: []string{
				"h1 not a header",
			},
		},
		{
			name: "disallows leading whitespace",
			in: []string{
				" h1.not a header",
			},
			out: []string{
				" h1.not a header",
			},
		},
		{
			name: "h0 -> h1",
			in: []string{
				"h0.My important header",
			},
			out: []string{
				"#My important header",
			},
		},
		{
			name: "h1 -> h2",
			in: []string{
				"h1.My important header",
			},
			out: []string{
				"##My important header",
			},
		},
		{
			// TODO: This is questionable, should it be changed?
			name: "h6 -> h7",
			in: []string{
				"h6.My important header",
			},
			out: []string{
				"#######My important header",
			},
		},
		{
			name: "h7 is not a thing",
			in: []string{
				"h7.not a header",
			},
			out: []string{
				"h7.not a header",
			},
		},
	},
	"code": []test{
		{
			name: "supports language marker",
			in: []string{
				"{code:go}",
				`func main() { fmt.Println("hello world") }`,
				"{code}",
			},
			out: []string{
				"```go",
				`func main() { fmt.Println("hello world") }`,
				"```",
			},
		},
	},
	"bold": []test{
		{
			name: "basic",
			in: []string{
				"supports *bold*",
			},
			out: []string{
				"supports **bold**",
			},
		},
	},
	"italic": []test{
		{
			name: "basic",
			in: []string{
				"supports _italic_",
			},
			out: []string{
				"supports *italic*",
			},
		},
	},
	"monospace": []test{
		{
			name: "basic",
			in: []string{
				"supports {{monospace}}",
			},
			out: []string{
				"supports `monospace`",
			},
		},
	},
	"ins": []test{
		{
			name: "basic",
			in: []string{
				"supports +inserts+",
			},
			out: []string{
				// TODO: Should it really keep the +?
				"supports <ins>+inserts+</ins>",
			},
		},
	},
	"sup": []test{
		{
			name: "basic",
			in: []string{
				"supports ^superscript^",
			},
			out: []string{
				"supports <sup>superscript</sup>",
			},
		},
	},
	"sub": []test{
		{
			name: "basic",
			in: []string{
				"supports ~subscript~",
			},
			out: []string{
				"supports <sub>subscript</sub>",
			},
		},
	},
	"strikethrough": []test{
		{
			name: "basic",
			in: []string{
				"supports -strikethrough- text",
			},
			out: []string{
				"supports ~~strikethrough~~ text",
			},
		},
	},
	"noformat": []test{
		{
			name: "basic",
			in: []string{
				"{noformat} supports noformat text",
			},
			out: []string{
				"``` supports noformat text",
			},
		},
	},
	"unnamedlinks": []test{
		{
			name: "basic",
			in: []string{
				"supports [unnamed links]",
			},
			out: []string{
				"supports <unnamed links>",
			},
		},
	},
	"images": []test{
		{
			name: "basic",
			in: []string{
				"supports !images.jpeg!",
			},
			out: []string{
				"supports ![](images.jpeg)",
			},
		},
	},
	"namedlinks": []test{
		{
			name: "basic",
			in: []string{
				"supports [named links|https://example.com]",
			},
			out: []string{
				"supports [named links](https://example.com)",
			},
		},
	},
	"blockquote": []test{
		{
			name: "disallows leading whitespace",
			in: []string{
				" bq.not a quote",
			},
			out: []string{
				" bq.not a quote",
			},
		},
		{
			name: "requires space after period",
			in: []string{
				"bq.not a quote",
			},
			out: []string{
				"bq.not a quote",
			},
		},
		{
			name: "simple positive case",
			in: []string{
				"bq. fourscore and seven",
			},
			out: []string{
				"> fourscore and seven",
			},
		},
	},
	"color": []test{
		{
			name: "strips color",
			in: []string{
				"{color:royalblue}The color of this text is royalblue.{color}",
			},
			out: []string{
				"The color of this text is royalblue.",
			},
		},
	},
	"th": []test{
		// TODO: fix
		{
			name: "basic",
			in: []string{
				" || col 1 || col 2 || col 3 ||",
			},
			out: []string{
				"| col 1 | col 2 | col 3 |",
				"|---|---|---|",
			},
		},
	},
}

var mixedFeatureTests = []test{
	{
		name: "unordered and bold compatibility",
		in: []string{
			"* *bold* unordered",
		},
		out: []string{
			"  * **bold** unordered",
		},
	},
}

func TestUnorderedList(t *testing.T) {
	for _, test := range singleFeatureTests["unordered list"] {
		t.Run("ul:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestOrderedList(t *testing.T) {
	for _, test := range singleFeatureTests["ordered list"] {
		t.Run("ol:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestHeader(t *testing.T) {
	for _, test := range singleFeatureTests["header"] {
		t.Run("header:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestCode(t *testing.T) {
	// TODO: Double check and fix whitespace expectations.
	for _, test := range singleFeatureTests["code"] {
		t.Run("code:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestBold(t *testing.T) {
	for _, test := range singleFeatureTests["bold"] {
		t.Run("bold:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestItalic(t *testing.T) {
	for _, test := range singleFeatureTests["italic"] {
		t.Run("italic:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestMonospace(t *testing.T) {
	for _, test := range singleFeatureTests["monospace"] {
		t.Run("monospace:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestIns(t *testing.T) {
	for _, test := range singleFeatureTests["ins"] {
		t.Run("ins:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestSup(t *testing.T) {
	for _, test := range singleFeatureTests["sup"] {
		t.Run("sup:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestSub(t *testing.T) {
	for _, test := range singleFeatureTests["sub"] {
		t.Run("sub:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestStrikethrough(t *testing.T) {
	for _, test := range singleFeatureTests["strikethrough"] {
		t.Run("strikethrough:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestNoformat(t *testing.T) {
	for _, test := range singleFeatureTests["noformat"] {
		t.Run("noformat:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestUnnamedLinks(t *testing.T) {
	for _, test := range singleFeatureTests["unnamedlinks"] {
		t.Run("unnamedlinks:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestImages(t *testing.T) {
	for _, test := range singleFeatureTests["images"] {
		t.Run("images:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestNamedLinks(t *testing.T) {
	for _, test := range singleFeatureTests["namedlinks"] {
		t.Run("namedlinks:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestBlockQuote(t *testing.T) {
	for _, test := range singleFeatureTests["blockquote"] {
		t.Run("blockquote:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestColor(t *testing.T) {
	for _, test := range singleFeatureTests["color"] {
		t.Run("color:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestTh(t *testing.T) {
	for _, test := range singleFeatureTests["th"] {
		t.Run("th:"+test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}

func TestMixedFeatures(t *testing.T) {
	for _, test := range mixedFeatureTests {
		t.Run(test.name, func(tt *testing.T) {
			actual := ToMarkdown(test.in)
			if !reflect.DeepEqual(actual, test.out) {
				t.Errorf("Expected %+v, but got %+v", test.out, actual)
			}
		})
	}
}
