package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	//并不需要使用其API，只需要执行该包的init方法（加载MySQL是驱动程序）
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Torrent struct {
	InfoHash string        `json:"infohash"`
	FileName string        `json:"name"`
	Files    []interface{} `json:"files"`
	Length   int           `json:"length"`
}

var (
	db *sqlx.DB
)

func main() {
	db, _ = sqlx.Open("mysql", "root:11521@tcp(localhost:3306)/dht")
	defer db.Close()

	// 读取文件
	fi, err := os.Open("C:\\Users\\snow\\go\\src\\github.com\\shiyanhui\\dht\\sample\\spider\\dht.log")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReaderSize(fi, 1024*1024*8)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		//fmt.Println(string(a))
		torrent := Torrent{}
		if len(a) > 0 {
			if err := json.Unmarshal(a, &torrent); err == nil {
				//fmt.Println(torrent)
				insert(torrent)
			} else {
				fmt.Println(err)
			}
		}
	}

}

func insert(t Torrent) {
	str, _ := json.Marshal(t.Files)
	_, e := db.Exec("replace into torrent(info_hash, file_name, files, length) values(?, ?, ?, ?);", t.InfoHash, t.FileName, string(str), t.Length)
	if e != nil {
		fmt.Println("err=", e, t.InfoHash)
		return
	}
}
