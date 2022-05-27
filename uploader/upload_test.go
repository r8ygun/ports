package uploader

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"ports/reader"
	"testing"
)

type MockRepository struct {
	NumCalls int
}

func (mr *MockRepository) Write(ctx context.Context, result reader.ReadResult) error {
	mr.NumCalls++
	return nil
}

func TestUpload(t *testing.T) {
	t.Run("writes the correct number of ports", func(t *testing.T) {
		portCount := 10
		readerCh := make(chan reader.ReadResult, portCount)
		for i := 0; i < portCount; i++ {
			readerCh <- reader.ReadResult{
				Id: fmt.Sprintf("ID-%d", i),
				Port: &reader.Port{
					Code: fmt.Sprintf("CODE-%d", i),
				},
			}
		}
		close(readerCh)

		repo := &MockRepository{}
		logBuffer := &bytes.Buffer{}
		logger := log.New(logBuffer, "", 0)

		uploader := Uploader{
			Reader:     readerCh,
			Repository: repo,
			Logger:     logger,
		}

		uploadedPortCount := uploader.Upload(context.Background())

		if uploadedPortCount != portCount {
			t.Errorf("expected [%d] ports to be written, but got %d", portCount, uploadedPortCount)
		}
	})
}
