package controller

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func getUploadsDir() (uploads string) {
	exe, err := os.Executable() //获取当前exe所在目录
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exe)                //获取dir
	uploads = filepath.Join(dir, "uploads") //dir后面+ uploads
	return
}

// 将网络路径 :path 变成本地绝对路径
func UploadsController(c *gin.Context) {
	if path := c.Param("path"); path != "" {
		target := filepath.Join(getUploadsDir(), path)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+path)
		c.Header("Content-Type", "application/octet-stream")
		c.File(target) //将文件写到body里
	} else {
		c.Status(http.StatusNotFound)
	}
}
