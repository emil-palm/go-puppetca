package pson

import (
	"encoding/json"

	"golang.org/x/text/encoding/charmap"
)

type String string

func (s *String) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)

	if err != nil {
		return err
	}

	decoder := charmap.ISO8859_1.NewDecoder()
	str, err = decoder.String(str)

	if err == nil {
		*s = String(str)
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

func (s String) String() string {
	return string(s)
}
