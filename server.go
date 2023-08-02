package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// Server は、http.Serverのラッパー構造体。
type Server struct {
	*http.Server
	l net.Listener
}

// NewServer は、Server構造体の新しいポインタを返す。
func NewServer(l net.Listener, mux http.Handler) *Server {
	s := &Server{
		Server: &http.Server{Handler: mux},
		l:      l,
	}
	return s
}

// Run は、HTTPサーバーを起動する。
func (s *Server) Run(ctx context.Context) error {
	// リクエストの処理中に外部からの終了通知（SIGTERM）を受け取っても、
	// 処理が完了してからシャットダウンするようにする: Graceful Shutdown.
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)
	// 1. 別のgoroutineでHTTPサーバーを起動する
	eg.Go(func() error {
		// 4. シャットダウンによりサーバーが終了する
		// シャットダウンの場合、errはhttp.ErrServerClosedになる
		if err := s.Serve(s.l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			// 6a. run()の戻り値としてエラーを返す
			return err
		}
		// 6b. run()の戻り値としてnilを返す
		return nil
	})

	// 2. チャネルからの終了通知を待機する
	<-ctx.Done()
	// 3. 終了通知を受け取ったら、サーバーをシャットダウンさせる
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	// 5. サーバーが終了したことにより、待機していたeg.Wait()が終了する
	return eg.Wait()
}
