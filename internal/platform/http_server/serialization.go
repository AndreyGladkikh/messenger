package http_server

import (
	"bytes"
	"encoding/json"
)

func encode(v any) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// func decode[T any]() {
// 	json.NewDecoder()
// }
