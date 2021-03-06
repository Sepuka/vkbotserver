package api

import "net/http"

const (
	defaultOutput = `ok`
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func Response() []byte {
	return []byte(defaultOutput)
}
