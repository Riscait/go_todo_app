package config

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	// 環境変数に、「TOTO_ENV=dev」をセットする
	wantEnv := "dev"
	t.Setenv("TODO_ENV", wantEnv)
	// 環境変数に、「ポート番号 = 3333」をセットする
	wantPort := 3333
	t.Setenv("PORT", fmt.Sprint(wantPort))
	// 環境変数をマッピングした構造体を取得する
	got, err := New()
	if err != nil {
		t.Fatalf("cannot create config: %v", err)
	}
	// ENV設定が環境変数にセットした値と同じか確認する
	if got.Env != wantEnv {
		t.Errorf("want %s, but got %s", wantEnv, got.Env)
	}
	// ポート番号が環境変数にセットした値と同じか確認する
	if got.Port != wantPort {
		t.Errorf("want %d, but got %d", wantPort, got.Port)
	}
}
