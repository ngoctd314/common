package ghttp

import (
	"encoding/json"
	"net/http"
)

// Unmarshal read body and Unmarshal it into T
// after read body, close it
func Unmarshal[T any](res *http.Response) (*T, error) {
	var result T

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return &result, nil
}
