package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

func handleChanList(w http.ResponseWriter, req *http.Request) {
	rows, err := db.Query(`
        SELECT name, address, sigkey, deckey
        FROM keys2
        WHERE enabled=1
        ORDER BY name ASC
        LIMIT 10000;
    `)
	if err != nil {
		http.Error(w, "Internal server error", 500)
		return
	}
	defer rows.Close()
	type Key2 struct {
		Name    string
		Address string
		Sigkey  string
		Deckey  string
	}
	type Info struct {
		Channels []Key2
		Count    int
	}
	var list []Key2
	var count int
	for rows.Next() {
		var name, address, sigkey, deckey string
		err := rows.Scan(&name, &address, &sigkey, &deckey)
		if err != nil {
			http.Error(w, "Internal server error", 500)
			return
		}
		if len(address) < 3 || address[:3] != "BM-" {
			address = "BM-" + address
		}
		list = append(list, Key2{
			Name:    name,
			Address: address,
			Sigkey:  sigkey,
			Deckey:  deckey,
		})
		count++
	}
	info := Info{
		Channels: list,
		Count:    count,
	}
	if err := tmpl.ExecuteTemplate(w, "chanlistpage", info); err != nil {
		logrus.Error(err)
	}
}
