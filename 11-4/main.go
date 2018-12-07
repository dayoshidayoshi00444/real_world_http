package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"io"
	"net/http"
	"os"
	"strings"
)

var clientID = "9febfdeb2d3ff325ff7a"
var clientSecret = "5fe81408ca573b11c89cd3c557b1c849fa80fae0"
var redirectURL = "https://localhost:18888"
var state = "https://github.com/dayoshidayoshi00444/dotfiles"

func main() {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user:email", "gist"},
		Endpoint:     github.Endpoint,
	}
	// これをこれから初期化する
	var token *oauth2.Token

	file, err := os.Open("access_token.json")
	if os.IsNotExist(err) {
		// 初回アクセス
		// まず認可画面のURLを取得
		url := conf.AuthCodeURL(state, oauth2.AccessTypeOnline)

		// コールバックを受け取るウェブサーバーをセットアップ
		code := make(chan string)
		var server *http.Server
		server = &http.Server{
			Addr: ":18888",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// クエリパラメータからcodeを取得し，ブラウザを閉じる
				w.Header().Set("Content-Type", "text/html")
				io.WriteString(w, "<html><script>window.open('about:blank','_self').close()</script></html>")
				w.(http.Flusher).Flush()
				code <- r.URL.Query().Get("code")
				// サーバーも閉じる
				server.Shutdown(context.Background())
			}),
		}
		go server.ListenAndServe()

		// ブラウザで認可画面を開く
		// Githubの認可が完了すれば上記のサーバーにリダイレクトされて，Handlerが実行される
		open.Start(url)

		// 取得したコードをアクセストークンに交換
		token, err = conf.Exchange(oauth2.NoContext, <-code)
		if err != nil {
			panic(err)
		}

		// アクセストークンをファイルに保存
		file, err := os.Create("access_token.json")
		if err != nil {
			panic(err)
		}
		json.NewEncoder(file).Encode(token)
	} else if err == nil {
		// 一度認可してローカルに保存ずみ
		token = &oauth2.Token{}
		json.NewDecoder(file).Decode(token)
	} else {
		panic(err)
	}
	client := oauth2.NewClient(oauth2.NoContext, conf.TokenSource(oauth2.NoContext, token))

	// gist投稿
	gist := `{
		"description": "API example",
		"public": true,
		"files": {
			"hello_from_rest_api.txt": {
				"content": "Hello World"
			}
		}
	}`

	// 投稿
	resp2, err := client.Post("https://api.github.com/gists", "application/json", strings.NewReader(gist))
	if err != nil {
		panic(err)
	}
	fmt.Println(resp2.Status)
	defer resp2.Body.Close()

	// 結果をパースする
	type GistResult struct {
		Url string `json:"html_url"`
	}
	gistResult := &GistResult{}
	err = json.NewDecoder(resp2.Body).Decode(&gistResult)
	if err != nil {
		panic(err)
	}
	if gistResult.Url != "" {
		// ブラウザで開く
		open.Start(gistResult.Url)
	}

}
