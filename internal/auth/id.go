package auth

import (
	"encoding/json"
)

type ID string

func (n *ID) UnmarshalJSON(b []byte) error {
	var result string

	// Add quotes to JSON to parse 123 as "123"
	if rune(b[0]) != '"' {
		b = []byte(`"` + string(b) + `"`)
	}

	err := json.Unmarshal(b, &result)
	*n = ID(result)
	return err
}
