package handlers

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"regexp"
	"strconv"
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
	Nsfw     int    `json:"nsfw" db:"nsfw"`
}

// List 列出指定目录下所有文件
func ListByKeyword(keyword string, page string) []Torrent {
	db, _ := sqlx.Open("sqlite3", "_torrent.db")
	defer db.Close()

	var torrents []Torrent
	// 关键字
	re := regexp.MustCompile(`[\s|.]+`)
	keyword = re.ReplaceAllString(keyword, "%")
	// page 为空判断
	pageNo := 1
	if page != "" {
		pageNo, _ = strconv.Atoi(page)
	}

	e := db.Select(&torrents, "select * from torrent where file_name like ? limit "+strconv.Itoa((pageNo-1)*20)+",20;", "%"+keyword+"%")
	if e != nil {
		log.Println("Error: " + e.Error())
		return make([]Torrent, 0)
	}

	if torrents == nil {
		torrents = make([]Torrent, 0)
	}

	return torrents
}
