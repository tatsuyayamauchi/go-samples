package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"google.golang.org/grpc/test/bufconn"
)

func NewServer(l *bufconn.Listener, msg string) *http.Server {
	return &http.Server{
		Addr: l.Addr().String(),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, msg)
		}),
	}
}

func NewClient(l *bufconn.Listener) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return l.DialContext(ctx)
			},
		},
	}
}
