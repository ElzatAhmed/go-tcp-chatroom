package server

import (
	"strings"
)

func parseProtocol(bytes []byte) (string, string, error) {
	protocol := strings.Trim(string(bytes), "\n")
	parsed := strings.Split(protocol, ";")
	if len(parsed) < 2 {
		return "", "", errProtocolError
	}
	return parsed[0], parsed[1], nil
}
