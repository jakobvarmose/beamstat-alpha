package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

type Comment struct {
	Sender     string
	SenderName string
	IsReply    bool
	Subject    string
	Text       string
	SafeText   template.HTML
	Pending    bool
	Received   string
	Body       string
	IsExtended bool
}

type Thread2 struct {
	Subject  string
	Hash     string
	Last     string
	Count    string
	Comments []*Comment
}

type Channel2 struct {
	Name    string
	Address string
	Exists  bool
	Enabled bool
	Threads []*Thread2
}

type Info2 struct {
	Channel *Channel2
	Thread  *Thread2
}

func formatTime(t int64) string {
	if t == 9223372036854775807 {
		return "Pending"
	}
	now := time.Now().Unix()
	diff := now - t
	if diff < -12*60*60 {
		return "Pending"
	}
	if diff < 60 {
		return "Just now"
	} else if diff < 120 {
		return "1 minute ago"
	} else if diff < 3600 {
		return fmt.Sprintf("%d minutes ago", diff/60)
	} else {
		return time.Unix(t, 0).Format("Jan 2 15:04")
	}
}

func handleThread(w http.ResponseWriter, req *http.Request, name string, threadHash string) {
	switch name {
	case "UPPERCASE_GENERAL":
		name = "GENERAL"
	case "aneki", "aneki/":
		name = "/aneki/"
	case "b", "b/":
		name = "/b/"
	case "pol", "pol/":
		name = "/pol/"
	case "tech", "tech/":
		name = "/tech/"
	}
	comments1, err := d.Comments(threadHash)
	if err != nil {
		http.Error(w, "Internal server error", 500)
		return
	}
	var comments []*Comment
	for _, comment := range comments1 {
		comments = append(comments, &Comment{
			Sender:     comment.Sender,
			SenderName: comment.SenderName,
			IsReply:    comment.IsReply,
			Subject:    comment.Subject,
			Text:       comment.Text,
			SafeText:   safeText(comment.Text),
			Pending:    comment.Pending,
			Received:   formatTime(comment.Received),
			Body:       comment.Body,
			IsExtended: comment.IsExtended,
		})
	}
	var subject string
	if len(comments) > 0 {
		subject = comments[0].Subject
	}
	channel, err := getChan(name)
	if err != nil {
		http.Error(w, "Internal server error", 500)
		return
	}
	info := Info2{
		Channel: channel,
		Thread: &Thread2{
			Comments: comments,
			Subject:  subject,
		},
	}
	if err := tmpl.ExecuteTemplate(w, "thread", info); err != nil {
		logrus.Error(err)
	}
}
