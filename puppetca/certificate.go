package puppetca

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/camptocamp/go-puppetca/puppetca/model"
	"github.com/camptocamp/go-puppetca/puppetca/pson"
	"github.com/pkg/errors"
)

const (
	CertificateStateSigned    = "signed"
	CertificateStateRevoked   = "revoked"
	CertificateStateRequested = "requested"
)

type CertificateSave struct {
	DesiredState pson.String `json:"desired_state"`
	TTL          int         `json:"cert_ttl,omitempty"`
}

// https://www.puppet.com/docs/puppet/8/server/http_certificate

func (c *Client) DownloadCertificateNamed(name string) (string, error) {
	pem, err := c.Get(c.NewRequest().SetPath("/certificate/%s", name))
	if err != nil {
		return "", errors.Wrapf(err, "failed to retrieve certificate %s", name)
	}

	// Response should always be a string with the content of the certificate
	return string(pem), nil
}

// https://www.puppet.com/docs/puppet/8/server/http_certificate_status#find

func (c *Client) GetCertificateNamed(name string) (*model.Certificate, error) {

	data, err := c.Get(c.NewRequest().SetPath("/certificate_status/%s", name))
	if err != nil {
		return nil, err
	}
	var cert model.Certificate
	err = json.Unmarshal(data, &cert)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

func (c *Client) GetCertificate(certificate model.Certificate) (*model.Certificate, error) {
	return c.GetCertificateNamed(string(certificate.Name))
}

// https://www.puppet.com/docs/puppet/8/server/http_certificate_status#search

func (c *Client) ListCertificates(state string) ([]*model.Certificate, error) {
	req := c.NewRequest().SetPath("/certificate_statuses/:any_key")
	if state != "" {
		req = req.AddQueryString("state", state)
	}

	data, err := c.Get(req)
	if err != nil {
		return nil, err
	}

	certificates := make([]model.Certificate, 0)
	err = json.Unmarshal(data, &certificates)
	if err != nil {
		return nil, err
	}

	return certificates, nil
}

// https://www.puppet.com/docs/puppet/8/server/http_certificate_status#save

func (c *Client) SaveCertificateNamed(name string, update CertificateSave) error {
	req := c.NewRequest().SetPath("certificate_statuses/%s", name).SetJSONBody(update)
	_, err := c.Put(req)
	return err
}

func (c *Client) SaveCertificate(cert model.Certificate, update CertificateSave) error {
	return c.SaveCertificateNamed(string(cert.Name), update)
}

func (c *Client) RevokeCertificateNamed(name string) error {
	return c.SaveCertificateNamed(name, CertificateSave{DesiredState: CertificateStateRevoked})
}

func (c *Client) RevokeCertificate(certificate model.Certificate) error {
	return c.RevokeCertificateNamed(string(certificate.Name))
}

/*
Cause the certificate authority to discard all SSL information regarding a host (including any certificates, certificate requests, and keys).
This does not revoke the certificate if one is present.

https://www.puppet.com/docs/puppet/8/server/http_certificate_status#delete
*/
func (c *Client) DeleteCertificateNamed(name string) error {
	resp, err := c.Delete(c.NewRequest().SetPath("/certificate_status/%s", name))
	if err != nil {
		return err
	}

	if string(resp) == "Nothing was deleted" {
		return fmt.Errorf("No certificate was deleted named %s", name)
	}

	return nil

}

func (c *Client) DeleteCertificate(cert *model.Certificate) error {
	return c.DeleteCertificateNamed(string(cert.Name))
}

// https://www.puppet.com/docs/puppet/8/server/http_certificate_clean

func (c *Client) CleanCertificateNames(certificates ...string) (cleaned, skipped []string, err error) {
	body := struct {
		CertNames []string `json:"certnames"`
	}{
		CertNames: certificates,
	}
	resp, err := c.Put(c.NewRequest().SetPath("/clean").SetJSONBody(body))

	if err != nil {
		return nil, nil, err
	}

	skipped_string := strings.Replace(string(resp), "The following certs do not exist and cannot be revoked: ", "", -1)
	var skipped_certs []string
	json.Unmarshal([]byte(skipped_string), &skipped_certs)

	cleaned = make([]string, 0)
	for _, cert := range certificates {
		for _, skipped := range skipped_certs {
			if cert == skipped {
				goto NEXT
			}
		}
		cleaned = append(cleaned, cert)
	NEXT:
	}

	return cleaned, skipped, err
}

func (c *Client) CleanCertificates(certificates ...model.Certificate) (cleaned, skipped []string, err error) {
	names := make([]string, len(certificates))
	for idx, cert := range certificates {
		names[idx] = string(cert.Name)
	}

	return c.CleanCertificateNames(names...)
}
