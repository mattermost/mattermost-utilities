package jira

import (
	"fmt"
	"testing"
)

// TODO: Implement some tests for real.
func TestAll(t *testing.T) {
	text := []string{
		"		*** 3 indent unordered",
		"# 1 indent ordered",
		"h1 not a header",
		" h1.not a header",
		"h1.My important header",
		"h6.Am I an h6?",
		"{code:go}",
		`func main() { fmt.Println("hello world") }`,
		"{code}",
		"supports *bold*, _italic_, and {{monospace}}",
		"supports +inserts+",
		"supports ^superscript^",
		"supports ~subscript~",
		"supports -strikethrough- text",
		"{noformat} supports noformat text",
		"supports [unnamed links]",
		"supports !images.jpeg!",
		"supports [named links|https://example.com]",
		" bq.not a quote",
		"bq.not a quote",
		"bq. fourscore and seven",
		"{color:royalblue}The color of this text is royalblue.{color}",
		" || col 1 || col 2 || col 3 ||",
	}

	for _, line := range ToMarkdown(text) {
		fmt.Printf("%q\n", line)
	}

	t.Errorf("need to implement proper tests")
}
