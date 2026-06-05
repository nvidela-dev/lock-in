package tui

import (
	"fmt"
	"strings"
)

func parseConfirmation(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "y":
		return true, nil
	case "n":
		return false, nil
	default:
		return false, fmt.Errorf("type y or n")
	}
}
