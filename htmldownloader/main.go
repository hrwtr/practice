package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/PuerkitoBio/goquery"
)

// リクエスト
type Request struct {
	url   string
	depth int
}

// 結果
type Result struct {
	err error
	url string
}

// チャンネル
type Channels struct {
	req  chan Request
	res  chan Result
	quit chan int
}

// チャンネルを取得。
func NewChannels() *Channels {
	return &Channels{
		req:  make(chan Request, 10),
		res:  make(chan Result, 10),
		quit: make(chan int, 10),
	}
}

// 指定された URL の Web ページを取得し、ページに含まれる URL の一覧を取得。
func Fetch(u string) (urls []string, err error) {
	baseUrl, err := url.Parse(u)
	if err != nil {
		return
	}

	resp, err := http.Get(baseUrl.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}

	urls = make([]string, 0)
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			reqUrl, err := baseUrl.Parse(href)
			if err == nil {
				urls = append(urls, reqUrl.String())
			}
		}
	})

	return
}

// クロール。
func Crawl(url string, depth int, ch *Channels) {
	defer func() { ch.quit <- 0 }()

	// WebページからURLを取得
	urls, err := Fetch(url)

	// 結果送信
	ch.res <- Result{
		url: url,
		err: err,
	}

	if err == nil {
		for _, url := range urls {
			// 新しいリクエスト送信
			ch.req <- Request{
				url:   url,
				depth: depth - 1,
			}
		}
	}
}

// クロールの深さの初期値
const crawlerDepthDefault = 2

// クロールの深さ
var crawlerDepth int

func main() {
	flag.IntVar(&crawlerDepth, "depth", crawlerDepthDefault, "クロールする深さを指定。")
	flag.Parse()

	i := 0

	if len(flag.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "URLを指定してください。")
		os.Exit(1)
	}
	startUrl := flag.Arg(0)
	if crawlerDepth < 1 {
		crawlerDepth = crawlerDepthDefault
	}

	chs := NewChannels()
	urlMap := make(map[string]bool)

	// 最初のリクエスト
	chs.req <- Request{
		url:   startUrl,
		depth: crawlerDepth,
	}

	fmt.Println("URL is ", startUrl, ", depth is ", crawlerDepth)
	//time.Sleep(7 * time.Second) // 7秒待つ
	// ワーカーの数
	wc := 0

	done := false
	for !done {
		select {
		case res := <-chs.res:
			if res.err == nil {
				fmt.Printf("Success %s\n", res.url)
				set(res.url, i)
				i++
			} else {
				fmt.Fprintf(os.Stderr, "Error %s\n%v\n", res.url, res.err)
			}
		case req := <-chs.req:
			if req.depth == 0 {
				break
			}

			if urlMap[req.url] {
				// 取得済み
				break
			}
			urlMap[req.url] = true

			wc++
			go Crawl(req.url, req.depth, chs)
		case <-chs.quit:
			wc--
			if wc == 0 {
				done = true
			}
		}
	}
}

func set(htmlurl string, i int) {
	// pUrl := flag.String("url", htmlurl, "URL to be processed")
	// fmt.Println("Download URL is ", htmlurl)
	// flag.Parse()
	// url := *pUrl
	// if url == "" {
	// 	fmt.Fprintf(os.Stderr, "Error: empty URL!\n")
	// 	return
	// }

	filename := path.Base(htmlurl)
	fmt.Println("Checking if " + filename + " exists ...")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		download(htmlurl, filename, i)
		fmt.Println(filename + " saved!")
	} else {
		fmt.Println(filename + " already exists!")
	}
}

func download(htmlurl, filename string, i int) {

	fmt.Println("Downloading " + htmlurl + " ...")
	resp, err := http.Get(htmlurl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//filename = fmt.Sprintf("%03d", i) + "_" + filename
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	io.Copy(f, resp.Body)
}
