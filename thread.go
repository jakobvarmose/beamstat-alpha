package main

import (
	"fmt"
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

	rows, err := db.Query(`
		select  coalesce(
			(select name from addresses where address=sender), sender
		), comment, received,
		sender, body, pending, subject, extended
        from channels
        where thread_hash = ? and received > unix_timestamp() - 3*28*24*60*60
        order by received asc
        limit 1000
	`, threadHash)
	//parent_hash, comment_hash, hash,
	var comments []*Comment
	if err != nil {
		logrus.Errorln(err.Error())
		http.Error(w, "Internal server error", 500)
		return
	}
	var subject string
	for rows.Next() {
		var comment Comment
		var received int64
		err := rows.Scan(
			&comment.SenderName,
			&comment.Text,
			&received,
			&comment.Sender,
			&comment.Body,
			&comment.Pending,
			&subject,
			&comment.IsExtended,
		)
		if err != nil {
			logrus.Errorln(err.Error())
		}
		comment.Received = formatTime(received)
		comments = append(comments, &comment)
	}
	rows.Close()
	channel, err := getChan(name)
	if err != nil {
		logrus.Errorln(err.Error())
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
