package reader

import (
	"encoding/json"
	"fmt"
	"io"
)

//Read will attempt to read and parse JSON Port objects from the given io.Reader.
//Parsed objects will be written to the returned unbuffered channel.
//Read will close the returned channel when there are no more objects to write.
func Read(r io.Reader) <-chan ReadResult {
	results := make(chan ReadResult)
	go func() {
		defer close(results)

		jsonDecoder := json.NewDecoder(r)

		_, err := jsonDecoder.Token()
		if err != nil {
			results <- ReadResult{
				Error: err,
			}
			return
		}

		for jsonDecoder.More() {
			name, err := jsonDecoder.Token()
			if err != nil {
				results <- ReadResult{
					Error: err,
				}
				return
			}

			var p Port
			err = jsonDecoder.Decode(&p)
			if err != nil {
				results <- ReadResult{
					Error: fmt.Errorf("error reading %s - %w", name, err),
				}
			}

			results <- ReadResult{
				Id:   name.(string),
				Port: &p,
			}
		}

		_, err = jsonDecoder.Token()
		if err != nil {
			results <- ReadResult{
				Error: err,
			}
		}

	}()
	return results
}
