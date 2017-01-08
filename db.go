package main

import (
	"github.com/ziutek/mymysql/autorc"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

func Open(user string, pass string, name string) *autorc.Conn {
	db := autorc.New("tcp", "", "127.0.0.1:3306", user, pass, name)
	db.Register("set names utf8")

	return db
}

var stmts = make(map[string]*autorc.Stmt)

func Query(db *autorc.Conn, query string, arguments ...interface{}) []mysql.Row {
	_, ok := stmts[query]
	if !ok {
		stmt2, err := db.Prepare(query)
		if err != nil {
			panic(err)
		}
		stmts[query] = stmt2
	}
	stmt := stmts[query]
	rows, _, err := stmt.Exec(arguments...)
	if err != nil {
		panic(err)
	}
	return rows
}
