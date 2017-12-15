package utils

import (
	"errors"
	"net"
)

func GetCurrentIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", errors.New("Could not determine current IP; there was an error listing network interfaces")
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("Could not determine current IP; it this the current only has a loopback IP ?")
}
