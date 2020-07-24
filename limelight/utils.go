package limelight

import (
	"fmt"
	"strings"
)

func splitSeparatedTriple(str, separator string) (string, string, string, error) {
	s := strings.SplitN(str, separator, 3)

	if len(s) != 3 || s[0] == "" || s[1] == "" || s[2] == "" {
		return "", "", "", fmt.Errorf("unexpected format (expected '<string>%s<string>%s<string>'): %s", separator, separator, str)
	}

	return s[0], s[1], s[2], nil
}

func splitSeparatedPair(str, separator string) (string, string, error) {
	s := strings.SplitN(str, separator, 2)

	if len(s) != 2 || s[0] == "" || s[1] == "" {
		return "", "", fmt.Errorf("unexpected format (expected '<string>%s<string>'): %s", separator, str)
	}

	return s[0], s[1], nil
}
