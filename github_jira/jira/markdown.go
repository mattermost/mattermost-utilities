package jira

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var reOl = regexp.MustCompile(`^[ \t]*(\*+)\s+`)

func ol(jira []string) {
	for i := range jira {
		line := jira[i]
		matches := reOl.FindSubmatchIndex([]byte(line))
		if len(matches) >= 4 {
			numStars := matches[3] - matches[2]
			replacement := strings.Repeat("  ", numStars) + "* "
			jira[i] = replacement + line[matches[1]:]
		}
	}
}

var reUl = regexp.MustCompile(`^[ \t]*(#+)\s+`)

func ul(jira []string) {
	for i := range jira {
		line := jira[i]
		matches := reUl.FindSubmatchIndex([]byte(line))
		if len(matches) >= 4 {
			numHashtag := matches[3] - matches[2]
			replacement := strings.Repeat("  ", numHashtag) + "1. "
			jira[i] = replacement + line[matches[1]:]
		}
	}
}

var reH = regexp.MustCompile(`^h([0-6])\.(.*)$`)

func h(jira []string) {
	for i := range jira {
		line := jira[i]
		matches := reH.FindStringSubmatch(line)
		if len(matches) >= 3 {
			headerNumberStr := matches[1]
			headerNumber, err := strconv.Atoi(headerNumberStr)
			if err != nil {
				log.Fatalf("Expected header tag to be a number, but got %q: %+v\n", headerNumberStr, err)
			}

			// TODO: Code this is ported from returns 1-7.
			// Should this be capped at 6 as h7 does not exist?
			jira[i] = strings.Repeat("#", headerNumber+1) + matches[2]
		}
	}
}

var reTh = regexp.MustCompile(`^[ \t]*((?:\|\|.*?)+\|\|)[ \t]*$`)
var reHeaderSeparator = regexp.MustCompile(`\|\|`)
var reHeaderMarkers = regexp.MustCompile(`\|[^|]+`)

// TODO: This function does not yet work.
// There seem to be incompatibilities between the regex it uses
// and golang not having the same concepts of multiline regexes.
// When entering the regex regexr.com, if the multiline option is checked
// with the text:
//
// 355
// adsffad
// ||col 1||col 2||col 3||
// || col 1 || col 2 || col 3 ||
// || col 1 || col 2 || col 3 ||
// ||
// asdf
//
// then matches are found
func th(jira []string) []string {
	jiraString := strings.Join(jira, "\n")
	groups := reTh.FindStringSubmatch(jiraString)
	// fmt.Printf("has groups?\n")
	// fmt.Printf("%+v\n", groups)
	if len(groups) < 2 {
		return jira
	}
	// fmt.Printf("has groups yay\n")
	singleBarred := reHeaderSeparator.ReplaceAllString(groups[1], "|")
	// fmt.Printf("singleBarred: %q\n", singleBarred)
	headerMarkers := reHeaderMarkers.ReplaceAllString(singleBarred, "| --- ")

	return strings.Split("\n"+singleBarred+"\n"+headerMarkers, "\n")
}

var reCode = regexp.MustCompile(`\{code(:([a-z]+))?([:|]?(title|borderStyle|borderColor|borderWidth|bgColor|titleBGColor)=.+?)*\}((.|\n)*?)\{code\}`)

func code(jira []string) []string {
	jiraString := strings.Join(jira, "\n")
	withCodeBlocks := reCode.ReplaceAllString(jiraString, "```$2\n$5\n```")

	return strings.Split(withCodeBlocks, "\n")
}

var reBold = regexp.MustCompile(`\*(\S.*)\*`)
var reItalic = regexp.MustCompile(`\_(\S.*)\_`)
var reMonospace = regexp.MustCompile(`\{\{([^}]+)\}\}`)

// Citations (buggy): Copied from old version and not converted over
// result = re.sub(r'\?\?((?:.[^?]|[^?].)+)\?\?', r'<cite>\1</cite>', result)

var reInsert = regexp.MustCompile(`\+([^+]*)\+`)
var reSuperscript = regexp.MustCompile(`\^([^^]*)\^`)
var reSubscript = regexp.MustCompile(`~([^~]*)~`)
var reStrikethrough = regexp.MustCompile(`(\s+)-(\S+.*?\S)-(\s+)`)
var rePreformatted = regexp.MustCompile(`{noformat}`)
var reUnnamedLink = regexp.MustCompile(`\[([^|]+)\]`)
var reImage = regexp.MustCompile(`!(.+)!`)
var reNamedLink = regexp.MustCompile(`\[(.+?)\|(.+)\]`)
var reSingleParagraphBlockquote = regexp.MustCompile(`^bq\.\s+`)
var reRemoveColor = regexp.MustCompile(`\{color:[^}]+\}(.*)\{color\}`)
var rePanelToTable = regexp.MustCompile(`\{panel:title=([^}]*)\}\n?(.*?)\n?\{panel\}`)
var reTrimTable = regexp.MustCompile(`^[ \t]*\|`)

func panelToTable(jira []string) []string {
	jiraString := strings.Join(jira, "\n")
	asTable := rePanelToTable.ReplaceAllString(jiraString, "\n| $1 |\n| --- |\n| $2 |")

	fmt.Printf("%q\n", asTable)
	return strings.Split(asTable, "\n")
}

func ToMarkdown(input []string) []string {
	tmp := input
	ol(tmp)
	ul(tmp)
	h(tmp)
	tmp = th(tmp)
	tmp = code(tmp)
	for i := range tmp {
		tmp[i] = reBold.ReplaceAllString(tmp[i], "**$1**")
		tmp[i] = reItalic.ReplaceAllString(tmp[i], "*$1*")
		tmp[i] = reMonospace.ReplaceAllString(tmp[i], "`$1`")
		// TODO: Should this do the whole match, or should it do the submatch,
		// like superscript and subscript are doing?
		tmp[i] = reInsert.ReplaceAllString(tmp[i], "<ins>$0</ins>")
		tmp[i] = reSuperscript.ReplaceAllString(tmp[i], "<sup>$1</sup>")
		tmp[i] = reSubscript.ReplaceAllString(tmp[i], "<sub>$1</sub>")
		tmp[i] = reStrikethrough.ReplaceAllString(tmp[i], "$1~~$2~~$3")
		tmp[i] = rePreformatted.ReplaceAllString(tmp[i], "```")
		tmp[i] = reUnnamedLink.ReplaceAllString(tmp[i], "<$1>")
		tmp[i] = reImage.ReplaceAllString(tmp[i], "![]($1)")
		tmp[i] = reNamedLink.ReplaceAllString(tmp[i], "[$1]($2)")
		tmp[i] = reSingleParagraphBlockquote.ReplaceAllString(tmp[i], "> ")
		tmp[i] = reRemoveColor.ReplaceAllString(tmp[i], "$1")
		tmp[i] = rePanelToTable.ReplaceAllString(tmp[i], "$1")
	}

	return tmp
}
