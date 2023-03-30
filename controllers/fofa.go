package controllers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type FofaController struct {
	beego.Controller
}

type QueryResult struct {
	Error      bool       `json:"error"`
	ConsumedFp int        `json:"consumed_fpoint"`
	Size       int        `json:"size"`
	Page       int        `json:"page"`
	Mode       string     `json:"mode"`
	Query      string     `json:"query"`
	Results    [][]string `json:"results"`
}

func (this *FofaController) Get() {
	this.TplName = "fofa.tpl"
}

func (this *FofaController) Post() {
	// 处理FOFA查询请求
	keyword := this.GetString("query")
	email := "962850765@qq.com"               // 替换为你的FOFA账户邮箱
	key := "af3c11358202ae64d889dab2b8caa559" // 替换为你的FOFA账户API Key
	encodedKeyword := base64.URLEncoding.EncodeToString([]byte(keyword))
	timeStamp := strconv.Itoa(int(time.Now().Unix()))
	sign := fmt.Sprintf("%s%d%s", encodedKeyword, timeStamp, key)
	signHmac := hmac.New(sha1.New, []byte(key))
	signHmac.Write([]byte(sign))
	signStr := base64.URLEncoding.EncodeToString(signHmac.Sum(nil))
	apiUrl := fmt.Sprintf("https://fofa.info/api/v1/search/all?email=%s&key=%s&qbase64=%s&page=1&size=1000&fields=ip,port,protocol,host,domain,os,server,icp,title", email, key, encodedKeyword)
	req, _ := http.NewRequest("GET", apiUrl, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", signStr)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var result QueryResult
	bodyString := string(body)
	if err := json.Unmarshal([]byte(bodyString), &result); err != nil {
		fmt.Println("解析 JSON 字符串出错：", err)
	}
	this.Data["Results"] = result.Results
	this.TplName = "fofa.tpl"
}
