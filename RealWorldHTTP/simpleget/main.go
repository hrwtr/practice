package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	//"net/http"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/idna"
)

// URL is default value
const URL = "http://localhost:18888"
const GIT = "http://github.com"

func main() {
	//getMethod()
	//getMethodWithQuery()
	//headMethod()
	//postMethod()
	//postFileMethod()
	//postObjectMethod()
	//postMultipartFormData()
	//postMime()
	//getCookie()
	//proxy()
	//localFileAccess()
	//deleteMethod()
	domainChange()

}

// curl http://localhost:18888
func getMethod() {
	// GET request を投げる
	// resp にはhttp.Request型のオブジェクトが入る
	resp, _ := http.Get(URL)

	// Responseのステータスを文字列で表示
	log.Println("Status:", resp.Status)

	// Responseのステータスコードで表示
	log.Println("StatusCode:", resp.StatusCode)

	// ResponseのHeaderを表示
	log.Println("Headers:", resp.Header)

	// ResponseのHeaderのContent-Lengthを表示
	log.Println("Header.Content-Length:", resp.Header.Get("Content-Length"))

	// close処理
	defer resp.Body.Close()

	// Responseのbodyをバイト列に変換
	body, _ := ioutil.ReadAll(resp.Body)

	// バイト列をstringに変換して出力
	log.Println(string(body))
}

// curl -G --data-urlencode "query=hello world" http://localhost:18888
func getMethodWithQuery() {

	values := url.Values{
		"query": {"hello world"},
	}

	resp, _ := http.Get(URL + "?" + values.Encode())
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
}

//	curl --head http://localhost:188888
func headMethod() {
	resp, _ := http.Head(URL)
	log.Println("Status:", resp.Status)
	log.Println("Status:", resp.Header)
}

//  curl -d test=value http://localhost:188888
func postMethod() {
	values := url.Values{
		"test": {"value"},
	}
	resp, _ := http.PostForm(URL, values)
	log.Println("Status", resp.Status)

}

//curl -T main.go -H "Content-Type: text/plain" http://localhost:188888
func postFileMethod() {
	file, _ := os.Open("main.go")
	resp, _ := http.Post(URL, "text/plain", file)
	log.Println("Status", resp.Status)
}

// throw object made in Program.
func postObjectMethod() {
	reader := strings.NewReader("テキスト")
	resp, _ := http.Post(URL, "text/plain", reader)
	log.Println("Status", resp.Status)
}

// curl -F "name=michael Jackson" -F "thumbnail=photo.jpg" http://localhost:188888
func postMultipartFormData() {
	// マルチパートをバイト列として格納するバッファを宣言
	var buffer bytes.Buffer

	// マルチパートを組み立てるwriterの作成
	writer := multipart.NewWriter(&buffer)

	//ファイル以外のフィールドはWriteField()で登録
	writer.WriteField("name", "Michael Jackson")

	//ファイル書き込みのio.Writer の作成
	fileWriter, _ := writer.CreateFormFile("thumbnail", "photo.jpg")

	// ファイルを開く
	readFile, _ := os.Open("photo.jpg")
	defer readFile.Close()

	//開いたファイルをio.Writerにコピーする
	io.Copy(fileWriter, readFile)
	writer.Close()

	// io.Writerの内容をバッファに書き込んでPostする
	resp, _ := http.Post(URL, writer.FormDataContentType(), &buffer)
	log.Println("Status", resp.Status)
}

// 任意のContent-Typeを指定できる
func postMime() {
	// マルチパートをバイト列として格納するバッファを宣言
	var buffer bytes.Buffer

	// マルチパートを組み立てるwriterの作成
	writer := multipart.NewWriter(&buffer)

	//ファイル以外のフィールドはWriteField()で登録
	writer.WriteField("name", "Michael Jackson")

	part := make(textproto.MIMEHeader)
	part.Set("Contet-type", "image/jpeg")
	part.Set("Content-Disposition", `form-data; name="thumbnail";filename="photo.jpg"`)
	fileWriter, _ := writer.CreatePart(part)
	readFile, _ := os.Open("photo.jpg")
	io.Copy(fileWriter, readFile)

	writer.Close()

	// io.Writerの内容をバッファに書き込んでPostする
	resp, _ := http.Post(URL, writer.FormDataContentType(), &buffer)
	log.Println("Status", resp.Status)
}

func getCookie() {
	//cookieを保存するインスタンスの作成
	jar, _ := cookiejar.New(nil)

	// cookieを保存可能なhttp.Clientインスタンスを作成
	client := http.Client{
		Jar: jar,
	}
	// cookieは1回目で保存、2回目で送信するので2回アクセスする
	for i := 0; i < 2; i++ {
		fmt.Println(i)
		//http.Getではなくclient.Getを使う
		resp, _ := client.Get(URL + "/cookie")
		dump, _ := httputil.DumpResponse(resp, true)
		log.Println(string(dump))
	}
}

// curl -x http://localhost:188888 http://github.com
func proxy() {
	proxyUrl, _ := url.Parse(URL)
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	resp, _ := client.Get(GIT)
	dump, _ := httputil.DumpResponse(resp, true)
	log.Println(string(dump))
}

// curl file:///Users/$USERNAME/Working/Develop/go/src/practice/RealWorldHTTP/simpleget/main.go
func localFileAccess() {
	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
	client := http.Client{
		Transport: transport,
	}
	resp, _ := client.Get("file://.main.go")
	dump, _ := httputil.DumpResponse(resp, true)
	log.Println(string(dump))
}

// curl -X DELETE http://localhost:188888
func deleteMethod() {

	client := &http.Client{}
	//http.Requestはhttp.NewRequest()というビルダーで生成する
	// 引数はメソッド、URL、ボディ
	request, _ := http.NewRequest("DELETE", URL, nil)
	resp, _ := client.Do(request)
	dump, _ := httputil.DumpResponse(resp, true)
	log.Println(string(dump))
}

// curl -H "Content-Type=@image/jpeg" -d "@image.jpeg" $URL

// ドメインチェンジ
func domainChange() {
	src := "握力王"
	ascii, _ := idna.ToASCII(src)
	fmt.Printf("%s -> %s\n", src, ascii)
}
