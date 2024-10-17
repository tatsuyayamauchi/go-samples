package server_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/tatsuyayamauchi/go-samples/util"
	server "github.com/tatsuyayamauchi/go-samples/virtual-conn-http-server"
	"google.golang.org/grpc/test/bufconn"
)

func TestParallelHTTPServer(t *testing.T) {
	t.Parallel()

	for i := range 100 {
		t.Run(fmt.Sprintf("parallel http server test(%d)", i), func(t *testing.T) {
			t.Parallel()

			l := bufconn.Listen(1024 * 1024)
			body := util.GenRandomString(20)

			svr := server.NewServer(l, body)
			go func() {
				if err := svr.Serve(l); err != nil && err != http.ErrServerClosed {
					t.Error(err)
				}
			}()
			defer svr.Shutdown(context.Background())

			client := server.NewClient(l)

			for range 100 {
				resp, err := client.Get(fmt.Sprintf("http://%s/", l.Addr().String()))
				if err != nil {
					t.Error(err)
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
				}

				actualBody, err := bufio.NewReader(resp.Body).ReadBytes('\n')
				if err != nil && err != io.EOF {
					t.Error(err)
				}
				if string(actualBody) != body {
					t.Errorf("body = %s, want %s", string(actualBody), body)
				}
			}
		})
	}
}
