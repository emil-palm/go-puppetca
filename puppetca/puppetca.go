package puppetca

import (
	"github.com/pkg/errors"
)

// Following functions are backported

// GetCertByName returns the certificate of a node by its name
func (c *Client) GetCertByName(nodename, env string) (string, error) {
	pem, err := c.DownloadCertificateNamed(nodename)
	if err != nil {
		return "", errors.Wrapf(err, "failed to retrieve certificate %s", nodename)
	}
	return pem, nil
}

// DeleteCertByName deletes the certificate of a given node
func (c *Client) DeleteCertByName(nodename, env string) error {
	err := c.DeleteCertificateNamed(nodename)
	if err != nil {
		return errors.Wrapf(err, "failed to delete certificate %s", nodename)
	}
	return nil
}

// GetRequest return a submitted CSR if any
func (c *Client) GetRequest(nodename, env string) (string, error) {

	pem, err := c.DownloadCertificateRequest(nodename)
	if err != nil {
		return "", errors.Wrapf(err, "failed to retrieve certificate %s", nodename)
	}
	return pem, nil
}

// SubmitRequest submits a CSR
func (c *Client) SubmitRequest(nodename, pem, env string) error {
	// Content-Type: text/plain
	err := c.CreateCertificateRequest(nodename, pem)
	if err != nil {
		return errors.Wrapf(err, "failed to submit CSR %s", nodename)
	}
	return nil
}

// SignRequest signs a CSR
func (c *Client) SignRequest(nodename, env string) error {
	err := c.SignCertificateRequestNamed(nodename, 0)
	if err != nil {
		return errors.Wrapf(err, "failed to sign CSR %s", nodename)
	}
	return nil
}
