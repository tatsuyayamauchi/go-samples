//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
)

func EventProvider() (*Event, error) {
	wire.Build(
		NewMessage,
		NewGreeter,
		NewEvent,
	)
	return nil, nil
}
