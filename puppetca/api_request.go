package puppetca

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type Request struct {
	baseurl string
	headers map[string]string
	query   url.Values
	path    string
	body    []byte
}

func NewRequest(baseurl string) *Request {
	return &Request{
		baseurl: baseurl,
		headers: make(map[string]string, 0),
		query:   make(url.Values, 0),
	}
}

func (r *Request) AddHeader(header, value string) *Request {
	r.headers[header] = value
	return r
}

func (r *Request) SetHeaders(headers map[string]string) *Request {
	r.headers = headers
	return r
}

func (r *Request) AddQueryString(key, value string) *Request {
	if value != "" {
		r.query.Add(key, value)
	}

	return r
}

func (r *Request) SetPath(format string, a ...interface{}) *Request {
	r.path = fmt.Sprintf(format, a...)
	return r
}

func (r *Request) SetBody(data []byte) *Request {
	r.body = data
	return r
}

func (r *Request) SetJSONBody(a interface{}) *Request {
	data, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	return r.SetBody(data)
}

func (r *Request) Build(method string) (*http.Request, error) {

	uri := fmt.Sprintf("%s%s", r.baseurl, r.path)
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create http request for URL %s", uri)
	}

	req.URL.RawQuery = r.query.Encode()
	for k, v := range r.headers {
		req.Header.Add(k, v)
	}

	if r.body != nil {
		req.Body = io.NopCloser(bytes.NewReader(r.body))
	}

	return req, nil
}
