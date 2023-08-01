package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNewMux は、NewMux関数で定義したルーティングが意図通りかのテストを行う。
func TestNewMux(t *testing.T) {
	// ServeHTTPに渡すためのモックをhttptestパッケージで作成する。
	// それぞれhttp.ResponseWriterと*http.Requestのモック
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	// sutはテスト対象システム (System Under Test)のこと。
	sut := NewMux()
	sut.ServeHTTP(w, r)

	// テスト対象ハンドラーを実行してレスポンスを取得する。
	resp := w.Result()
	t.Cleanup(func() { _ = resp.Body.Close() })
	// ステータスコードのチェック
	if resp.StatusCode != http.StatusOK {
		t.Error("want status code 200, but", resp.StatusCode)
	}
	// レスポンスボディを取得。
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	// レスポンスボディのチェック
	want := `{ "status": "ok" }`
	if string(got) != want {
		t.Errorf("want %q, but got %q:", want, got)
	}
}
