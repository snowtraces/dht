package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shiyanhui/dht/sample/web/repository"
	"net/http"
)

var (
	db *sqlx.DB
)

type BaseRes struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// ListTorrents 列出指定目录下所有文件
func ListTorrents(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	keyword := values.Get("keyword")
	page := values.Get("page")

	// 1. 查询数据库
	torrents := handlers.ListByKeyword(keyword, page)

	jsonBytes, _ := json.Marshal(&torrents)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonBytes)
}

func Error(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(401)
	fmt.Fprintln(w, "认证后才能访问该接口")
}

func exception(w http.ResponseWriter, r *http.Request, msg string) {
	res := BaseRes{}
	res.Code = 500
	res.Msg = msg
	jsonBytes, _ := json.Marshal(&res)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonBytes)
}
