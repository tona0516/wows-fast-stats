package domain

import (
	"errors"
	"slices"
)

// https://ja.wikipedia.org/wiki/ISO_3166-1
type CountryCode struct {
	ValueObject[string]
}

func NewCountryCode(value string) (CountryCode, error) {
	// 空文字はそれ以外の国
	if !slices.Contains([]string{"jp", "cn", "kr", ""}, value) {
		return CountryCode{}, errors.New("invalid_country_code: " + value)
	}

	return CountryCode{ValueObject[string]{value}}, nil
}
