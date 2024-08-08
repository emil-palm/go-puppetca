package puppetca

import (
	"encoding/json"

	"github.com/camptocamp/go-puppetca/puppetca/model"
)

// https://www.puppet.com/docs/puppet/8/server/http_certificate_sign

func (c *Client) DownloadCertificateRequest(name string) (string, error) {
	data, err := c.Get(c.NewRequest().SetPath("certificate_request/%s", name))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (c *Client) CreateCertificateRequest(name string, pem string) error {
	req := c.NewRequest().SetPath("certificate_request/%s", name)
	req.SetBody([]byte(pem))
	_, err := c.Put(req)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DestroyCertificateRequestNamed(name string) error {
	_, err := c.Delete(c.NewRequest().SetPath("certificate_request/%s", name))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DestroyCertificateRequest(certificate model.Certificate) error {
	return c.DestroyCertificateRequestNamed(string(certificate.Name))
}

// https://www.puppet.com/docs/puppet/8/server/http_certificate_status#search-unsigned-certs-csrs
func (c *Client) ListCertificateRequests() ([]model.Certificate, error) {
	return c.ListCertificates(CertificateStateRequested)
}

// https://www.puppet.com/docs/puppet/8/server/http_certificate_status#save

func (c *Client) SignCertificateRequestNamed(name string, ttl int) error {
	save := CertificateSave{DesiredState: CertificateStateSigned}
	if ttl != 0 {
		save.TTL = ttl
	}
	return c.SaveCertificateNamed(name, save)
}

func (c *Client) SignCertificateRequest(certificate model.Certificate, ttl int) error {
	return c.SignCertificateRequestNamed(string(certificate.Name), ttl)
}

// https://www.puppet.com/docs/puppet/8/server/http_certificate_sign

type bulkSignBody struct {
	CertNames []string `json:"certnames"`
}

type bulkResponseBody struct {
	Signed []string `json:"signed"`
	NoCSR  []string `json:"no-csr"`
	Errors []string `json:"signing-errors"`
}

func (c *Client) BulkSignCertificateRequestsNamed(certificates ...string) (signed, missing_csr, sign_errors []string, err error) {
	body := bulkSignBody{
		CertNames: certificates,
	}

	data, err := c.Post(c.NewRequest().SetPath("sign").SetJSONBody(body))

	if err != nil {
		return nil, nil, nil, err
	}

	var respBody bulkResponseBody

	err = json.Unmarshal(data, &respBody)
	if err != nil {
		return nil, nil, nil, err
	}

	return respBody.Signed, respBody.NoCSR, respBody.Errors, nil
}

func (c *Client) BulkSignCertificateRequests(certificates ...model.Certificate) (signed, missing_csr, sign_errors []string, err error) {
	names := make([]string, len(certificates))
	for idx, cert := range certificates {
		names[idx] = string(cert.Name)
	}

	return c.BulkSignCertificateRequestsNamed(names...)
}

func (c *Client) SignAllCertificateRequests() (signed, missing_csr, sign_errors []string, err error) {

	data, err := c.Post(c.NewRequest().SetPath("sign/all"))

	if err != nil {
		return nil, nil, nil, err
	}

	var respBody bulkResponseBody

	err = json.Unmarshal(data, &respBody)
	if err != nil {
		return nil, nil, nil, err
	}

	return respBody.Signed, respBody.NoCSR, respBody.Errors, nil
}
