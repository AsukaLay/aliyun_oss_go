package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	res, err := http.Get("http://192.168.31.142:9999")
	if err != nil {
		fmt.Println(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Print(string(body))
	requestPara := &Res{}
	err = json.Unmarshal(body, &requestPara)
	if err != nil {
		return
	}
	upLoadFileInfo := &UpLoadFileInfo{Path: "C:\\Users\\Administrator\\Desktop\\docker.json", Name: "docker"}
	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf) // body writer

	f, err := os.Open(upLoadFileInfo.Path)
	if err != nil {
		return
	}
	defer f.Close()

	p1w, _ := bw.CreateFormField("name")
	p1w.Write([]byte(upLoadFileInfo.Name))

	p2w, _ := bw.CreateFormField("key")
	p2w.Write([]byte(requestPara.Dir + "12315/docker.json"))

	p3w, _ := bw.CreateFormField("policy")
	p3w.Write([]byte(requestPara.Policy))

	p4w, _ := bw.CreateFormField("OSSAccessKeyId")
	p4w.Write([]byte(requestPara.Accessid))

	p5w, _ := bw.CreateFormField("success_action_status")
	p5w.Write([]byte(strconv.Itoa(200)))

	p6w, _ := bw.CreateFormField("callback")
	p6w.Write([]byte(requestPara.Callback))

	p7w, _ := bw.CreateFormField("signature")
	p7w.Write([]byte(requestPara.Signature))

	// file part1
	_, fileName := filepath.Split(upLoadFileInfo.Path)
	fw1, _ := bw.CreateFormFile("file", fileName)
	io.Copy(fw1, f)

	bw.Close() //write the tail boundry
	client := &http.Client{}
	reqa, err := http.NewRequest("POST", requestPara.Host, buf)
	fmt.Println(string(buf.String()))
	if err != nil {
		// handle error
		log.Fatal("生成请求失败！", err)
		return
	}
	reqa.Header.Add("Content-Type", bw.FormDataContentType())

	postResponse, err := client.Do(reqa)
	if err != nil {
		return
	}
	if err != nil {
		fmt.Println(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("关闭失败！", err)
		}
	}(postResponse.Body)
	postResponseBody, _ := ioutil.ReadAll(res.Body)
	fmt.Print(string(postResponseBody))

}

type UpLoadFileInfo struct {
	Path string
	Name string
}

type Res struct {
	Accessid  string `json:"accessid"`
	Host      string `json:"host"`
	Expire    int64  `json:"expire"`
	Policy    string `json:"policy"`
	Signature string `json:"signature"`
	Callback  string `json:"callback"`
	Dir       string `json:"dir"`
}


type Req struct {
	Name                  string `json:"name"`
	Key                   string `json:"key"`
	Policy                string `json:"policy"`
	OSSAccessKeyId        string `json:"OSSAccessKeyId"`
	Success_action_status int32  `json:"success_action_status"`
	Signature             string `json:"signature"`
	Callback              string `json:"callback"`
	FilePath              string `json:"file"`
}
