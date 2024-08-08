package model

import (
	"github.com/camptocamp/go-puppetca/puppetca/pson"
)

type Certificate struct {
	Name                 pson.String                 `json:"name"`
	DnsAltNames          []pson.String               `json:"dns_alt_names"`
	State                pson.String                 `json:"state"`
	AuthorizedExtensions map[pson.String]pson.String `json:"authorization_extensions"`
	SubjectAltNames      []pson.String               `json:"subject_alt_names"`
	NotBefore            pson.String                 `json:"not_before"`
	NotAfter             pson.String                 `json:"not_after"`
	Serial               int64                       `json:"serial"`
	Fingerprint          pson.String                 `json:"fingerprint"`
	Fingerprints         map[pson.String]pson.String `json:"fingerprints"`
}
