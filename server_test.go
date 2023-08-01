package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestServer_Run(t *testing.T) {
	// 利用可能なポート番号を動的に選択してリスナーを作成
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}
	// 1. キャンセル可能なContextオブジェクトを作る
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	// multiplexer: 複数の入力信号を1つの信号として出力する装置
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})
	// 2. 別goroutineでテスト対象であるrun関数を実行してHTTPサーバーを起動する
	eg.Go(func() error {
		s := NewServer(l, mux)
		return s.Run(ctx)
	})
	// 3. エンドポイントに対してGETリクエストを送信する
	path := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), path)
	// どんなポート番号でリッスンしているのか確認
	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()
	// レスポンスボディを読み込む（[]byte型）
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	// 6. HTTPサーバーの戻り値（文字列）を検証する
	want := fmt.Sprintf("Hello, %s!", path)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
	// 4. run関数に終了通知を送信する。
	cancel()
	// 5. Waitメソッド経由で、run関数の戻り値を検証する
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
