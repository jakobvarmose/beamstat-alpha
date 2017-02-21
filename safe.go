package main

import (
	"html"
	"html/template"
	"regexp"
)

var httpRegexp = regexp.MustCompile(`http(?:s?)://(?:[^\s">]+)`)

func safeText(str string) template.HTML {
	var out []byte
	index := 0
	loc := httpRegexp.FindStringIndex(str[index:])
	for loc != nil {
		out = append(out, html.EscapeString(str[index:index+loc[0]])...)
		out = append(out, `<a rel="nofollow" href="`...)
		out = append(out, html.EscapeString(str[index+loc[0]:index+loc[1]])...)
		out = append(out, `">`...)
		out = append(out, html.EscapeString(str[index+loc[0]:index+loc[1]])...)
		out = append(out, `</a>`...)
		index = index + loc[1]
		loc = httpRegexp.FindStringIndex(str[index:])
	}
	out = append(out, html.EscapeString(str[index:])...)
	return template.HTML(out)
}
