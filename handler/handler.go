package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrResponse はエラーレスポンスのJSONを表す。
type ErrResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// RespondJSON はレスポンスとしてJSONを返す。エラーが発生した場合はエラーレスポンスを返す。
func RespondJSON(
	ctx context.Context,
	w http.ResponseWriter,
	body any,
	statusCode int,
) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyJSON, err := json.Marshal(body)

	if err != nil {
		// JSONのエンコードに失敗した場合は500エラーを返す。
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		rsp := ErrResponse{
			Message: http.StatusText(code),
		}
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			fmt.Printf("write response error: %v", err)
		}
		return
	}
	// 正常系のステータスコードとレスポンスを返す。
	w.WriteHeader(statusCode)
	if _, err := fmt.Fprintf(w, "%s", bodyJSON); err != nil {
		fmt.Printf("write response error: %v", err)
	}
}
