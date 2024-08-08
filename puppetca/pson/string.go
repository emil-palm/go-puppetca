package pson

import (
	"encoding/json"

	"golang.org/x/text/encoding/charmap"
)

type String string

func (s *String) UnmarshalJSON(b []byte) error {
	decoder := charmap.ISO8859_1.NewDecoder()
	b, err := decoder.Bytes(b)

	if err == nil {
		*s = String(b)
	}

	return err
}

func (s *String) MarshalJSON() ([]byte, error) {
	encoder := charmap.ISO8859_1.NewEncoder()
	str, err := encoder.String(string(*s))
	if err != nil {
		return nil, err
	}

	return json.Marshal(str)
}
