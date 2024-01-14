package helpers

import (
	"regexp"
	"strconv"
)

func ValidatePort(port_str string, default_port ...int) (port int) {
	defaultPort := 1080
	if len(default_port) > 0 {
		defaultPort = default_port[0]
	}

	portRe := regexp.MustCompile(`^\d+$`)
	if !portRe.MatchString(port_str) {
		return defaultPort
	}

	port, err := strconv.Atoi(port_str)
	if err != nil {
		return defaultPort
	}

	if port < 1 || port > 65535 {
		return defaultPort
	}

	return port
}
