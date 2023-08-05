package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Riscait/go_todo_app/entity"
	"github.com/Riscait/go_todo_app/store"
	"github.com/Riscait/go_todo_app/testutil"
	"github.com/go-playground/validator/v10"
)

func TestAddTask_ServeHTTP(t *testing.T) {
	// このテストが他の並列テストと並列に実行されるようにする。
	t.Parallel()

	// 期待値の型定義。
	type want struct {
		// レスポンスのステータスコード。
		status int
		// レスポンスJSONのGoldenファイルのパス。
		rspFile string
	}

	// テストケースをmapで複数用意。
	tests := map[string]struct {
		// リクエストJSONのGoldenファイルのパス。
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_task/201_request.json.golden",
			want: want{
				status:  http.StatusCreated,
				rspFile: "testdata/add_task/201_response.json.golden",
			},
		},
		"validate error": {
			reqFile: "testdata/add_task/400_request.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/add_task/400_response.json.golden",
			},
		},
	}

	for n, tt := range tests {
		t.Run(n, func(t *testing.T) {
			t.Parallel()
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest(
			http.MethodPost,
			"/tasks",
			bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
		)

		sut := AddTask{
			Store: &store.TaskStore{
				Tasks: map[entity.TaskID]*entity.Task{},
			},
			Validator: validator.New(),
		}
		sut.ServeHTTP(w, r)

		resp := w.Result()
		testutil.AssertResponse(
			t,
			resp,
			tt.want.status,
			testutil.LoadFile(t, tt.want.rspFile),
		)
	}
}
