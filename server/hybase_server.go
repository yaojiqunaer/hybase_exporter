package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HyBaseItem struct {
	Unit  string `json:"unit"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Key   string `json:"key"`
}
type HyBaseNode struct {
	Id       string       `json:"id"`
	Ip       string       `json:"ip"`
	Sysname  string       `json:"sysname"`
	Itemlist []HyBaseItem `json:"itemlist"`
}

type Result struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    []HyBaseNode `json:"data"`
}

func genHyBaseNode() HyBaseNode {
	return HyBaseNode{
		Id:      "hybase1",
		Ip:      "192.168.210.57",
		Sysname: "hybase",
		Itemlist: []HyBaseItem{
			{
				"",
				"节点状态",
				"on",
				"server_status",
			},
			{
				"%",
				"CPU使用率",
				"0.17",
				"cpu_usage",
			},
			{
				"%",
				"磁盘使用率",
				"51.6",
				"disk_usage",
			},
			{
				"GB",
				"磁盘总空间",
				"1023.5",
				"disk_size",
			},
			{
				"GB",
				"磁盘剩余空间",
				"495.33",
				"disk_free",
			},
			{
				"MB/s",
				"磁盘读速率",
				"0.0",
				"disk_io_out",
			},
			{
				"MB/s",
				"磁盘写速率",
				"0.0",
				"disk_io_in",
			},
			{
				"MB/s",
				"网卡下行速率",
				"0.0",
				"nic_byte_in",
			},
			{
				"MB/s",
				"网卡上行速率",
				"0.0",
				"nic_byte_out",
			},
			{
				"次/s",
				"检索效率",
				"0.0",
				"search_rate",
			},
		},
	}
}

func main() {
	//海贝返回信息
	hybaseInfo := genHyBaseNode()
	var info []HyBaseNode
	info = append(info, hybaseInfo)
	result := Result{
		Code:    0,
		Message: "成功",
		Data:    info,
	}
	// 1.创建路由
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/public/status.do", func(c *gin.Context) {
		c.JSON(http.StatusOK, result)
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	err := r.Run(":5555")
	if err != nil {
		return
	}
}
