package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/shiyanhui/dht"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
)

type file struct {
	Path   []interface{} `json:"path"`
	Length int           `json:"length"`
}

type bitTorrent struct {
	InfoHash string `json:"infohash"`
	Name     string `json:"name"`
	Files    []file `json:"files,omitempty"`
	Length   int    `json:"length,omitempty"`
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func main() {
	var (
		fileName = "./dht.log"
		logFile  *os.File
		logErr   error
	)

	func() {
		//文件是否存在
		if Exists(fileName) {
			//使用追加模式打开文件
			logFile, logErr = os.OpenFile(fileName, os.O_APPEND, 0666)
			if logErr != nil {
				fmt.Println("打开文件错误：", logErr)
				return
			}
		} else {
			//不存在创建文件
			logFile, logErr = os.Create(fileName)
			if logErr != nil {
				fmt.Println("创建失败", logErr)
				return
			}
		}
	}()
	defer logFile.Close()

	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	w := dht.NewWire(65536, 1024, 256)
	go func() {
		for resp := range w.Response() {
			metadata, err := dht.Decode(resp.MetadataInfo)
			if err != nil {
				continue
			}
			info := metadata.(map[string]interface{})

			if _, ok := info["name"]; !ok {
				continue
			}

			bt := bitTorrent{
				InfoHash: hex.EncodeToString(resp.InfoHash),
				Name:     info["name"].(string),
			}

			if v, ok := info["files"]; ok {
				files := v.([]interface{})
				bt.Files = make([]file, len(files))

				for i, item := range files {
					f := item.(map[string]interface{})
					bt.Files[i] = file{
						Path:   f["path"].([]interface{}),
						Length: f["length"].(int),
					}
				}
			} else if _, ok := info["length"]; ok {
				bt.Length = info["length"].(int)
			}

			data, err := json.Marshal(bt)
			if err == nil {
				fmt.Printf("%s\n\n", data)
				io.WriteString(logFile, string(data)+"\n")
			}
		}
	}()
	go w.Run()

	config := dht.NewCrawlConfig()
	config.OnAnnouncePeer = func(infoHash, ip string, port int) {
		w.Request([]byte(infoHash), ip, port)
	}
	d := dht.New(config)

	d.Run()
}
