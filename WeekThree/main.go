package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	server := &http.Server{Addr: "127.0.0.1:8000", Handler: mux}

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return errors.Wrap(server.ListenAndServe(),"http服务报错")
	})
	
	g.Go(func() error {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		select {
		case <-signalChan:
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			return errors.Wrap(server.Shutdown(ctx),"http服务优雅退出")
		case <-ctx.Done():
			return errors.New("ctx done")
		}
	})

	err := g.Wait()
	if err != nil {
		fmt.Println(err)
		return
	}
}
