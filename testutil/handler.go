package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// AssertJSON はJSONの比較検査を行う。差分があれば差分を出力し、テストを失敗させる。
func AssertJSON(t *testing.T, want, got []byte) {
	// テストヘルパー関数としてマークし、複数のgoroutineから同時に呼び出せるようになる。
	t.Helper()

	var jw, jg any
	if err := json.Unmarshal(want, &jw); err != nil {
		t.Fatalf("cannot unmarshal want %q: %v", want, err)
	}
	if err := json.Unmarshal(got, &jg); err != nil {
		t.Fatalf("cannot unmarshal got %q: %v", got, err)
	}
	if diff := cmp.Diff(jg, jw); diff != "" {
		t.Errorf("JSON mismatch (-got +want):\n%s", diff)
	}
}

// AssertResponse はレスポンスの比較検査を行う。差分があれば差分を出力し、テストを失敗させる。
func AssertResponse(
	t *testing.T,
	got *http.Response,
	wStatusCode int,
	wBody []byte,
) {
	t.Helper()
	t.Cleanup(func() { _ = got.Body.Close() })

	gBody, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatal(err)
	}
	if got.StatusCode != wStatusCode {
		t.Errorf(
			"want status code %d, but got %d, body: %q",
			wStatusCode, got.StatusCode, gBody,
		)
	}
	if len(gBody) == 0 && len(wBody) == 0 {
		// 期待も実体も空のレスポンスボディなので、JSONの検証は不要。
		return
	}
	AssertJSON(t, wBody, gBody)
}

// LoadFile はファイルを読み込む。読み込みに失敗した場合はテストを失敗させる。
func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read from %q: %v", path, err)
	}
	return content
}
