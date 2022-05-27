package uploader

import (
	"context"
	"errors"
	"log"
	"ports/reader"
)

var (
	ErrUploadOperationCancelled = errors.New("the upload operation was cancelled")
)

type Uploader struct {
	Reader     <-chan reader.ReadResult
	Repository PortsRepositoryWriter
	Logger     *log.Logger
}

func New(reader <-chan reader.ReadResult, repo PortsRepositoryWriter, logger *log.Logger) *Uploader {
	return &Uploader{
		Reader:     reader,
		Repository: repo,
		Logger:     logger,
	}
}

//Upload reads Port objects and attempts to persist them.
//This method will return a count of the number of Ports it uploaded.
func (u *Uploader) Upload(ctx context.Context) int {
	uploaded := 0
	for {
		select {
		case result, more := <-u.Reader:
			if !more {
				return uploaded
			}
			if result.Error != nil {
				u.Logger.Println(result.Error)
				continue
			}
			err := u.Repository.Write(ctx, result)
			if err != nil {
				u.Logger.Println(err)
				continue
			}
			uploaded++
			u.Logger.Printf("Uploaded new document for [%s]\n", result.Id)
		case <-ctx.Done():
			u.Logger.Println(ErrUploadOperationCancelled)
			return uploaded
		}
	}
}
