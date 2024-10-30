package main

import (
	"context"
	"errors"
	"fmt"
	"goplus/mock_secware/pkg/config"
	"goplus/mock_secware/pkg/handlers"
	"net"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.OnSecwareTask)
	http.HandleFunc("/meta", handlers.OnGetMeta)
	http.HandleFunc("/health", handlers.OnGetHealth)

	ctx, cancelCtx := context.WithCancel(context.Background())
	newCtx, err := handlers.BuildContext(ctx)
	if err != nil {
		fmt.Print(err)
		return
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", config.PORT),
		BaseContext: func(l net.Listener) context.Context {
			return newCtx
		},
	}

	go func() {
		fmt.Printf("Listening on %s\n", server.Addr)
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
		} else if err != nil {
			fmt.Printf("error starting server: %s\n", err)
		}
		cancelCtx()
	}()

	<-ctx.Done()
}
