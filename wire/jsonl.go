package wire

import (
	"encoding/json"
)

func EncodeLine(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	data = append(data, '\n')
	return data, nil
}
