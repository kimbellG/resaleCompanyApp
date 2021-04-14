package models

import (
	"encoding/json"
	"net/http"
)

func DecodingAnswerForRegistration(r *http.Request) (*RegistrationInformation, error) {
	result := &RegistrationInformation{}

	if err := json.NewDecoder(r.Body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
