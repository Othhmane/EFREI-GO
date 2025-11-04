package utils

import (
	"bufio"
	"strconv"
	"strings"
)

// ReadLine lit une ligne depuis le reader et retire les espaces
func ReadLine(r *bufio.Reader) (string, error) {
	s, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(s), nil
}

// ParseInt convertit une chaîne en int
func ParseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}