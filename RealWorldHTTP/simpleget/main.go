package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

// URL is default value
const URL = "http://localhost:18888"

func main() {
	// GET request を投げる
	// resp にはhttp.Request型のオブジェクトが入る
	resp, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	// close処理
	defer resp.Body.Close()

	// Responseのbodyをバイト列に変換
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// バイト列をstringに変換して出力
	log.Println(string(body))
}
