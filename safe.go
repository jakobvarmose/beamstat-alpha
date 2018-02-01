package main

import (
	"html"
	"html/template"
	"regexp"
)

var httpRegexp = regexp.MustCompile(`http(?:s?)://(?:[^\s">]+)|\[chan\] \S+`)

func safeText(str string) template.HTML {
	var out []byte
	index := 0
	loc := httpRegexp.FindStringIndex(str[index:])
	for loc != nil {
		out = append(out, html.EscapeString(str[index:index+loc[0]])...)
		if str[index+loc[0]] == 'h' {
			out = append(out, `<a rel="nofollow" href="`...)
			out = append(out, html.EscapeString(str[index+loc[0]:index+loc[1]])...)
			out = append(out, `">`...)
		} else {
			out = append(out, `<a href="/chan/`...)
			out = append(out, html.EscapeString(str[index+loc[0]+7:index+loc[1]])...)
			out = append(out, `">`...)
		}
		out = append(out, html.EscapeString(str[index+loc[0]:index+loc[1]])...)
		out = append(out, `</a>`...)
		index = index + loc[1]
		loc = httpRegexp.FindStringIndex(str[index:])
	}
	out = append(out, html.EscapeString(str[index:])...)
	return template.HTML(out)
}
