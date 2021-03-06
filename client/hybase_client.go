// Package client
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"hybase_exporter/common/log"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	// hybasePublicStatus status url
	hybasePublicStatus = "/public/status.do"
)

// HyBaseResult HyBase Api基础返回实体
type HyBaseResult struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    []HyBaseNode `json:"data"`
}

// HyBaseItem 各个指标数据
type HyBaseItem struct {
	Unit  string `json:"unit"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Key   string `json:"key"`
}

// HyBaseNode HyBaseNode
type HyBaseNode struct {
	Id       string       `json:"id"`
	Ip       string       `json:"ip"`
	Sysname  string       `json:"sysname"`
	Itemlist []HyBaseItem `json:"itemlist"`
}

// HyBaseClient Get the hybase http client instance
type HyBaseClient struct {
	Protocol string
	Ip       string
	Port     int
}

// getHyBaseUrl Get the api address based on the underlying path and URI
func getHyBaseUrl(client HyBaseClient, uri string) string {
	return client.Protocol + "://" + client.Ip + ":" + strconv.Itoa(client.Port) + uri
}

var hbCli *http.Client = &http.Client{
	Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, time.Second*3)
			if err != nil {
				fmt.Println("dail timeout", err)
				return nil, err
			}
			return c, nil
		},
		IdleConnTimeout:       time.Second * 65,
		ResponseHeaderTimeout: time.Second * 3,
	},
}

// GetPublicStatus Send http request to hybase server and parse res to struct HyBaseResult
func GetPublicStatus(client HyBaseClient) HyBaseResult {
	url := getHyBaseUrl(client, hybasePublicStatus)
	log.Info.Log("msg", "send request to "+url)
	resp, err := hbCli.Get(url)
	if err != nil {
		log.Error.Log("msg", err)
		return HyBaseResult{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error.Log("msg", err)
		return HyBaseResult{}
	}
	return parseHyBaseResult(body)
}

// parseHyBaseResult parse hybase api result json bytes to struct
func parseHyBaseResult(result []byte) HyBaseResult {
	var hbr HyBaseResult
	err := json.Unmarshal(result, &hbr)
	if err != nil {
		log.Warn.Log("msg", "parse hybase json result error", "result", string(result), "err", err)
		return HyBaseResult{}
	}
	return hbr
}
