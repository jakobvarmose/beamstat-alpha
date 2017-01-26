package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/jakobvarmose/beamstat-alpha/dao"
)

func handleChanList(w http.ResponseWriter, req *http.Request) {
	keys, err := d.AllKeys()
	if err != nil {
		http.Error(w, "Internal server error", 500)
		return
	}
	type Info struct {
		Channels []*dao.Key
		Count    int
	}
	info := Info{
		Channels: keys,
		Count:    len(keys),
	}
	if err := tmpl.ExecuteTemplate(w, "chanlistpage", info); err != nil {
		logrus.Error(err)
	}
}
