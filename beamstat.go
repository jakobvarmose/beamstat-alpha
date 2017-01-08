package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

var tmpl *template.Template
var db *sql.DB

func init2() {
	dir := "templates"
	files, err := AssetDir(dir)
	if err != nil {
		panic(err)
	}
	tmpl = template.New("")
	for _, file := range files {
		_, err := tmpl.Parse(string(MustAsset(path.Join(dir, file))))
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	buf, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	config := make(map[string]string)
	if err := json.Unmarshal(buf, &config); err != nil {
		panic(err)
	}
	db, err = sql.Open(config["sql_driver"], config["sql_database"])
	if err != nil {
		panic(err)
	}
	/*	for _, row := range db.MustQuery(`select data from objects2 join objects on objects.hash=objects2.hash where type=0 limit 5`) {
		obj, err := types.UnmarshalObject(row.Bin(0))
		if err != nil {
			panic(err)
		}
		switch obj.Type {
		case 0:
			getpubkey, err := types.UnmarshalGetpubkey(obj.Payload)
			println("getpubkey", obj.Version, getpubkey.Tag, getpubkey.Ripe, err)
		case 1:
			println("pubkey")
		case 2:
			println("msg")
		}
	}*/
	_, err = db.Query("select 1")
	if err != nil {
		panic(err)
	}
	//fileServer := http.FileServer(http.Dir("/home/i/go/src/jakobvarmose/beamstat/static"))
	//http.Handle("/css/bootstrap.min.css", fileServer)
	init2()

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/css/bootstrap.min.css" {
			w.Header().Add("Content-Type", "text/css")
			w.Write(MustAsset("static/css/bootstrap.min.css"))
			return
		}
		if req.URL.Path == "/css/main.css" {
			w.Header().Add("Content-Type", "text/css")
			w.Write(MustAsset("static/css/main.css"))
			return
		}
		if req.URL.Path == "/" {
			handleIndex(w, req)
			return
		}
		if req.URL.Path == "/downloads" {
			handleDownload(w, req)
			return
		}
		match := regexp.MustCompile(`^/chan/(.+)/([0-9a-f]{64})$`).FindStringSubmatch(req.URL.Path)
		if match != nil {
			name, err := url.QueryUnescape(match[1])
			if err == nil {
				hash := match[2]
				handleThread(w, req, name, hash)
				return
			}
		}
		match = regexp.MustCompile(`^/chan/(.+)$`).FindStringSubmatch(req.URL.Path)
		if match != nil {
			name, err := url.QueryUnescape(match[1])
			if err == nil {
				handleChan(w, req, name)
				return
			}
		}
		/*		if
				 else if ok, _ := regexp.MatchString(`^/chan/[^/]+$`, req.URL.Path); ok {
					handleChan(w, req)
				} else if ok, _ := regexp.MatchString(`^/obj$`, req.URL.Path); ok {
					//handleObjList(w, req)
					http.NotFound(w, req)
				} else if ok, _ := regexp.MatchString(`^/obj/[0-9a-f]{64}$`, req.URL.Path); ok {
					//handleObj(w, req)
					http.NotFound(w, req)
				} else {*/
		http.NotFound(w, req)
	})
	logrus.Info("Listening on http://127.0.0.1:8002/chan/general")
	if err := http.ListenAndServe(":8002", nil); err != nil {
		logrus.Errorln(err.Error())
	}
}
