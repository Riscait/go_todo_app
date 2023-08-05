package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Riscait/go_todo_app/entity"
	"github.com/Riscait/go_todo_app/store"
	"github.com/go-playground/validator/v10"
)

type AddTask struct {
	Store     *store.TaskStore
	Validator *validator.Validate
}

func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// リクエストボディをマッピングする構造体を用意。
	var b struct {
		Title string `json:"title" validate:"required"`
	}
	// リクエストボディJSONをデコードする。失敗した場合は500エラーを返す。
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		if errors.Is(io.EOF, err) {
			RespondJSON(ctx, w, &ErrResponse{
				Message: "required request body",
			}, http.StatusBadRequest)
			return
		}
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// リクエストボディのバリデーションを行う。異常があれば400エラーを返す。
	if err := at.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	// DBに登録するタスク構造体を作成する。
	t := &entity.Task{
		Title:     b.Title,
		Status:    entity.TaskStatusTodo,
		CreatedAt: time.Now(),
	}
	// storeパッケージを使ってタスクを追加する。失敗した場合は500エラーを返す。
	id, err := at.Store.Add(t)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// 追加したタスクのIDをレスポンスとして返す。
	rsp := struct {
		ID int `json:"id"`
	}{ID: int(id)}
	RespondJSON(ctx, w, rsp, http.StatusCreated)
}
