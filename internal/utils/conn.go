package utils

import (
	"net"
	"time"
)

func IsOnline() bool {
	const googleDNS = "8.8.8.8:53"
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", googleDNS, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}