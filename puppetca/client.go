package puppetca

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func isFile(str string) bool {
	return strings.HasPrefix(str, "/")
}

// NewClient returns a new Client
func NewClient(baseURL, keyStr, certStr, caStr string) (c Client, err error) {
	// Load client cert
	var cert tls.Certificate
	if isFile(certStr) {
		if !isFile(keyStr) {
			err = fmt.Errorf("cert points to a file but key is a string")
			return
		}

		cert, err = tls.LoadX509KeyPair(certStr, keyStr)
		if err != nil {
			err = errors.Wrapf(err, "failed to load client cert from file %s", certStr)
			return c, err
		}
	} else {
		if isFile(keyStr) {
			err = fmt.Errorf("cert is a string but key points to a file")
			return c, err
		}

		cert, err = tls.X509KeyPair([]byte(certStr), []byte(keyStr))
		if err != nil {
			err = errors.Wrapf(err, "failed to load client cert from string")
			return c, err
		}
	}

	// Load CA cert
	var caCert []byte
	if isFile(caStr) {
		caCert, err = ioutil.ReadFile(caStr)
		if err != nil {
			err = errors.Wrapf(err, "failed to load CA cert at %s", caStr)
			return
		}
	} else {
		caCert = []byte(caStr)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{Transport: tr}
	c = Client{fmt.Sprintf("%s/%s", baseURL, "puppet-ca/v1"), httpClient}

	return
}

func (c *Client) NewRequest() *Request {
	return NewRequest(c.baseURL)
}

// Get performs a GET request
func (c *Client) Get(request *Request) ([]byte, error) {
	req, err := request.Build(http.MethodGet)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post performs a POST request
func (c *Client) Post(request *Request) ([]byte, error) {
	req, err := request.Build(http.MethodPost)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Put performs a PUT request
func (c *Client) Put(request *Request) ([]byte, error) {
	req, err := request.Build(http.MethodPut)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Delete performs a DELETE request
func (c *Client) Delete(request *Request) ([]byte, error) {
	req, err := request.Build(http.MethodDelete)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Do performs an HTTP request
func (c *Client) Do(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to %s URL %s", req.Method, req.URL)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("failed to %s URL %s, got: %s", req.Method, req.URL, resp.Status)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read body response")
	}

	return data, nil
}
