package loader

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RemoteFileLoader struct {
	Url string
}

func (r RemoteFileLoader) LoadBatch(batchSize int) (data []interface{}, atEof bool, err error) {
	if batchSize != 1 {
		return nil, false, errors.New("remote file loader must have a batch size of 1, at least for now")
	}

	// make the request
	resp, err := http.Get(r.Url)
	if err != nil {
		return nil, false, err
	}

	// enforce the need for a response body
	if resp.Body == nil {
		return nil, false, errors.New("empty response body")
	}

	// prepare a close on the response body
	defer resp.Body.Close()

	// check server response
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("bad status: %s", resp.Status)
	}

	// read the whole body into a byte array
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, false, err
	}

	// return the byte array as the first (and only) result
	result := make([]interface{}, 1)
	result[0] = b
	return result, true, nil
}
