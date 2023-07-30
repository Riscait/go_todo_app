package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func main() {
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}
	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", p, err)
	}
	if err := run(context.Background(), l); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, l net.Listener) error {
	s := &http.Server{
		// 引数で受け取ったnet.Listenerを使うので、Addrフィールドは指定しない
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	// 1. 別のgoroutineでHTTPサーバーを起動する
	eg.Go(func() error {
		// 4. シャットダウンによりサーバーが終了する
		// シャットダウンの場合、errはhttp.ErrServerClosedになる
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
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
