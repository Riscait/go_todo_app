package main

import (
	"net/http"

	"github.com/Riscait/go_todo_app/handler"
	"github.com/Riscait/go_todo_app/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// NewMux は、HTTPリクエストを処理するハンドラーを登録したServeMuxを返す。
// ルーティング定義を行っている。
// muxとは、複数の入力信号を1つの信号として出力するmultiplexerのこと。
// http.ServeMux型はルーティングの表現力が乏しい。
// 関数シグネチャを変更せずに実装できる`go-chi/chi`を使って改善した。
func NewMux() http.Handler {
	mux := chi.NewRouter()
	// HTTPサーバーが稼働中か確認するための/healthエンドポイント
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	v := validator.New()

	at := &handler.AddTask{Store: store.Tasks, Validator: v}
	mux.Post(
		"/tasks",
		at.ServeHTTP,
	)

	lt := &handler.ListTask{Store: store.Tasks}
	mux.Get(
		"/tasks",
		lt.ServeHTTP,
	)

	return mux
}
