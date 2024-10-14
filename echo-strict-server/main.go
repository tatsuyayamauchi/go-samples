package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/tatsuyayamauchi/go-samples/echo-strict-server/handler"
	"github.com/tatsuyayamauchi/go-samples/echo-strict-server/spec"
)

func main() {
	e := echo.New()

	spec.RegisterHandlers(e, spec.NewStrictHandler(handler.NewHandler(), nil))
	err := e.Start(":8080")
	if err != nil && err != http.ErrServerClosed {
		slog.Error(fmt.Sprintf("Start error: %v", err.Error()))
	}
}
