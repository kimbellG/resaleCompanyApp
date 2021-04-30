package models

import (
	"encoding/json"
	"net/http"
)

func DecodingPasswordAuthInfo(r *http.Request) (*PasswordAutheficationInfo, error) {
	result := &PasswordAutheficationInfo{}

	if err := json.NewDecoder(r.Body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func DecodingProvider(r *http.Request) (*Provider, error) {
	result := &Provider{}

	if err := json.NewDecoder(r.Body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
