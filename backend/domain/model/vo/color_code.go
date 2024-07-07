package vo

import (
	"errors"
	"regexp"
)

type ColorCode struct {
	ValueObject[string]
}

func NewColorCode(value string) (ColorCode, error) {
	regex := `^#(?:[0-9a-fA-F]{3}){1,2}$`
	isMatch, _ := regexp.MatchString(regex, value)
	if !isMatch {
		return ColorCode{}, errors.New("invalid_color_code: " + value)
	}

	return ColorCode{ValueObject[string]{value}}, nil
}
