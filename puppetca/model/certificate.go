package model

type Certificate struct {
	Name                 string            `json:"name"`
	DnsAltNames          []string          `json:"dns_alt_names"`
	State                string            `json:"state"`
	AuthorizedExtensions map[string]string `json:"authorization_extensions"`
	SubjectAltNames      []string          `json:"subject_alt_names"`
	NotBefore            string            `json:"not_before"`
	NotAfter             string            `json:"not_after"`
	Serial               int64             `json:"serial"`
	Fingerprint          string            `json:"fingerprint"`
	Fingerprints         map[string]string `json:"fingerprints"`
}
