package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type ChannelInfo struct {
	Index int
	Count int
	Name  string
	Last  string
}

type ActiveInfo struct {
	TypeString string
	Sub        bool
	Count      int
	Size       int
}

func getActiveList() []ActiveInfo {
	list := make([]ActiveInfo, 0)

	rows, err := db.Query(`
			select type, if(type=2 and version=1 and size=54, -1, version), count(*), sum(size) div 1000
			from objects
			where expires >= ?
			group by type, if(type=2 and version=1 and size=54, -1, version)
			limit 1000
		`, time.Now().Unix())
	if err != nil {
		panic(err)
	}
	var subs []ActiveInfo
	for rows.Next() {
		var t int
		var v int
		var count int
		var size int
		rows.Scan(&t, &v, &count, &size)
		var typeString string
		if t == 0 {
			typeString = "Getpubkey"
		} else if t == 1 {
			typeString = "Pubkey"
		} else if t == 2 {
			typeString = "Msg"
		} else if t == 3 {
			typeString = "Broadcast"
		} else {
			typeString = fmt.Sprintf("Type %d", t)
		}
		subs = append(subs, ActiveInfo{
			TypeString: fmt.Sprintf("%s v%s", typeString, v),
			Sub:        true,
			Count:      count,
			Size:       size,
		})
		/*TODO if i+1 == len(rows) || t != rows[i+1].Int(0) {
			list = append(list, ActiveInfo{
				TypeString: typeString,
				Sub:        false,
				Count:      count,
				Size:       size,
			})
			list = append(list, subs...)
			subs = nil
		}*/
	}
	rows.Close()
	return list
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	type Info struct {
		Channels []ChannelInfo
		Actives  []ActiveInfo
	}
	info := Info{
		getChanList(),
		getActiveList(),
	}
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
	tmpl.Execute(w, info)
}
