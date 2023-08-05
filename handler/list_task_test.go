package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Riscait/go_todo_app/entity"
	"github.com/Riscait/go_todo_app/store"
	"github.com/Riscait/go_todo_app/testutil"
)

func TestListTask_ServeHTTP(t *testing.T) {
	// 期待値の型定義。
	type want struct {
		// レスポンスのステータスコード。
		status int
		// レスポンスJSONのGoldenファイルのパス。
		rspFile string
	}

	// テストケースをmapで複数用意。
	tests := map[string]struct {
		tasks entity.Tasks
		want  want
	}{
		"ok": {
			tasks: []*entity.Task{
				{
					ID:     1,
					Title:  "Todo task title",
					Status: entity.TaskStatusTodo,
				},
				{
					ID:     2,
					Title:  "Doing task title",
					Status: entity.TaskStatusDoing,
				},
				{
					ID:     3,
					Title:  "Done task title",
					Status: entity.TaskStatusDone,
				},
			},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_task/200_response.json.golden",
			},
		},
		"empty": {
			tasks: []*entity.Task{},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_task/empty_response.json.golden",
			},
		},
	}

	for n, tt := range tests {
		// t.Runを使用したサブテストでは、goroutineが使われており、実行時にループ変数を参照する。
		// 最後のテストケースが繰り返されないように、このスコープでttを再定義している。
		tt := tt
		t.Run(n, func(t *testing.T) {
			// このテストが他の並列テストと並列に実行されるようにする。
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/tasks", nil)

			tasks := map[entity.TaskID]*entity.Task{}
			for _, t := range tt.tasks {
				tasks[t.ID] = t
			}
			// テスト対象のハンドラーを生成。
			sut := ListTask{
				// Storeには期待値となるタスクのリストをセット。
				Store: &store.TaskStore{
					Tasks: tasks,
				},
			}
			sut.ServeHTTP(w, r)

			rsp := w.Result()
			testutil.AssertResponse(
				t,
				rsp,
				tt.want.status,
				testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}
}
