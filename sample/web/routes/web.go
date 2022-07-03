package routes

import (
	"github.com/shiyanhui/dht/sample/web/handlers"
	"net/http"
)

// WebRoute 定义一个 WebRoute 结构体用于存放单个路由
type WebRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// WebRoutes 声明 WebRoutes 切片存放所有 Web 路由
type WebRoutes []WebRoute

// 定义所有 Web 路由
var webRoutes = WebRoutes{
	WebRoute{
		"ListTorrents",
		"GET",
		"/api/list_torrents",
		handlers.ListTorrents,
	},
}
