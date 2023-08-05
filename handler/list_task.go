package handler

import (
	"net/http"

	"github.com/Riscait/go_todo_app/entity"
	"github.com/Riscait/go_todo_app/store"
)

type ListTask struct {
	Store *store.TaskStore
}

// レスポンス用のタスク構造体。DB上のタスク構造体とは別。
type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// storeパッケージを使って全てのタスクを取得する。
	tasks := lt.Store.All()
	// タスクのスライスをレスポンスとして返す。※nilが返らないように{}で初期化をしておく。
	rsp := []task{}
	for _, t := range tasks {
		rsp = append(rsp, task{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		})
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
