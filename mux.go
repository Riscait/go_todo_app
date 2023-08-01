package main

import "net/http"

// NewMux は、HTTPリクエストを処理するハンドラーを登録したServeMuxを返す。
// ルーティング定義を行っている。
// muxとは、複数の入力信号を1つの信号として出力するmultiplexerのこと。
func NewMux() http.Handler {
	mux := http.NewServeMux()
	// HTTPサーバーが稼働中か確認するための/healthエンドポイント
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	return mux
}
