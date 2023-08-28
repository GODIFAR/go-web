package server

import (
	"embed"
	"go-web/server/controller"
	"go-web/server/ws"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var FS embed.FS //打包go的exe文件时，把frontend/dist/* 打包进去

func Run() {
	//设置gin模式
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	//打包的文件变成结构化目录staticFiles
	staticFiles, _ := fs.Sub(FS, "frontend/dist")

	router.POST("/api/v1/files", controller.FilesController)
	router.GET("/api/v1/qrcodes", controller.QrcodesController)
	router.GET("/uploads/:path", controller.UploadsController)
	router.GET("/api/v1/addresses", controller.AddressesController)
	router.POST("/api/v1/texts", controller.TextsController)
	hub := ws.NewHub()
	go hub.Run()
	router.GET("/ws", func(c *gin.Context) {
		ws.HttpController(c, hub)
	})
	router.StaticFS("/static", http.FS(staticFiles))
	//如果没有路由
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path               //用户的访问路径
		if strings.HasPrefix(path, "/static/") { //如果是静态文件
			reader, err := staticFiles.Open("index.html") //打开index.html
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()
			stat, err := reader.Stat() //读取stat(描述文件或目录的信息)
			if err != nil {
				log.Fatal(err)
			}
			//DataFromReader将指定的阅读器写入正文流并更新HTTP代码
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html", reader, nil)
		} else { //动态文件
			c.Status(http.StatusNotFound) //404
		}
	})
	router.Run(":27149")
}
