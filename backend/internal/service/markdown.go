package service

import (
	"bytes"
	"regexp"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var mdParser = goldmark.New(
	goldmark.WithExtensions(extension.GFM, extension.Table),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	goldmark.WithRendererOptions(html.WithHardWraps(), html.WithUnsafe()),
)

var markdownSymbols = regexp.MustCompile("[#*_\\[\\]()~`>|!\\-=]+")

func RenderMarkdown(source string) string {
	var buf bytes.Buffer
	if err := mdParser.Convert([]byte(source), &buf); err != nil {
		return source
	}
	return buf.String()
}

func WordCount(source string) int {
	plain := markdownSymbols.ReplaceAllString(source, " ")
	count := 0
	inWord := false
	for _, r := range plain {
		if unicode.Is(unicode.Han, r) {
			count++
			inWord = false
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if !inWord {
				count++
				inWord = true
			}
		} else {
			inWord = false
		}
	}
	return count
}

func ReadingMinutes(wordCount int) int {
	mins := wordCount / 300
	if mins < 1 {
		return 1
	}
	return mins
}
