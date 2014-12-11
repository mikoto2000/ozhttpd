package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// 現在時刻をフォーマットして返却
func now() string {
	return time.Now().Format("2006/01/02 15:04:05.999")
}

// ファイルを読み込む
func readFile(uri string) (content []byte, err error) {
	contents, err := ioutil.ReadFile(uri)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return nil, err
	}

	return contents, nil
}

// 適当に Content-Type を設定する
func getContentType(uri string) (contentType string) {
	switch true {
	case strings.HasSuffix(uri, ".htm"):
		fallthrough
	case strings.HasSuffix(uri, ".html"):
		contentType = "text/html"
	case strings.HasSuffix(uri, ".css"):
		contentType = "text/css"
	case strings.HasSuffix(uri, ".json"):
		fallthrough
	case strings.HasSuffix(uri, ".js"):
		contentType = "text/javascript"
	case strings.HasSuffix(uri, ".txt"):
		contentType = "text/plain"
	default:
		contentType = "application/octet-stream"
	}

	return
}

// サーバを表す構造体
type httpd struct {
	host    string
	port    uint
	docroot string
}

// リクエストを受け取った時の処理
func (httpd httpd) ServeHTTP(
	w http.ResponseWriter, r *http.Request) {

	// ログに出力
	now := now()
	uri := r.RequestURI
	fmt.Printf("[%v] %v\n", now, uri)

	// ファイルを探してレスポンスとして返却
	content, err := readFile(httpd.docroot + uri)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(404)
		fmt.Fprint(w, "404 File Not Found.")
	}

	// Content-Type の取得と設定
	contentType := getContentType(uri)
	w.Header().Set("Content-Type", contentType)

	binary.Write(w, binary.BigEndian, content)
}

func main() {

	// docroot のデフォルトを作成
	cd, _ := path.Split(os.Args[0])
	defaultDocroot := cd + "docroot"

	// 引数解析
	host := flag.String("h", "localhost", "listen host")
	port := flag.Uint("p", 8080, "listen port")
	docroot := flag.String("d", defaultDocroot, "docroot path")
	flag.Parse()

	// 設定適用
	httpd := httpd{*host, *port, *docroot}

	// 設定表示
	fmt.Printf("Start httpd.\n[Host: %v, Port: %v, Docroot: %v]\n",
		*host,
		*port,
		*docroot,
	)

	// httpd 実行
	http.Handle("/", httpd)
	http.ListenAndServe(fmt.Sprintf("%v:%v", *host, *port), nil)
}
