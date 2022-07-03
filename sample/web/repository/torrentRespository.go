package handlers

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type BaseRes struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Torrent struct {
	InfoHash string `json:"info_hash" db:"info_hash"`
	FileName string `json:"file_name" db:"file_name"`
	Files    string `json:"files" db:"files"`
	Length   int    `json:"length" db:"length"`
}

// List 列出指定目录下所有文件
func ListByKeyword(keyword string) []Torrent {
	db, _ := sqlx.Open("mysql", "root:11521@tcp(localhost:3306)/dht")
	defer db.Close()

	var torrents []Torrent
	e := db.Select(&torrents, "select * from torrent where file_name like ? limit 100;", "%"+keyword+"%")
	if e != nil {
		return nil
	}

	return torrents
}
