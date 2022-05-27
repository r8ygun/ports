package uploader

import (
	"context"
	"ports/reader"
)

type PortsRepositoryWriter interface {
	Write(ctx context.Context, result reader.ReadResult) error
}
