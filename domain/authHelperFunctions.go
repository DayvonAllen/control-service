package domain

import (
	"fmt"
	"strings"
)

func ExtractData(token string) ([]string, error) {
	if token == "" {
		return nil, fmt.Errorf("no token provided")
	}
	xs := strings.Split(token, " ")

	if len(xs) != 2 {
		return nil, fmt.Errorf("invalid token provided")
	}

	tokenValue := strings.Split(xs[1], "|")
	return tokenValue, nil
}
