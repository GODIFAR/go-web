package controller

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddressesController(c *gin.Context) {
	addrs, _ := net.InterfaceAddrs() //获取电脑在各个局域网的IP地址
	var result []string
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址 check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = append(result, ipnet.IP.String())
			}
		}
	}
	//转为json写入HTTP响应
	c.JSON(http.StatusOK, gin.H{"addresses": result})
}
