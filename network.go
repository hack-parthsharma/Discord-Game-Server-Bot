package network

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/pkg/errors"
)

type Network struct {
	IP string `json:"ip,omitempty"`
}

var endpoints = map[string]bool{
	"https://api.ipify.org?format=text": true,
	"https://myexternalip.com/text":     true,
	"https://v4.ident.me/":              true,
}

func GetPublicIP() (string, error) {
	for endpoint, enabled := range endpoints {
		if enabled {
			resp, err := http.Get(endpoint)
			if err != nil {
				endpoints[endpoint] = false
				return "", errors.Wrapf(err, "error getting IP address from endpoint: %v", endpoint)
			}

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()

			if err != nil {
				return "", errors.Wrapf(err, "unable to read response body, endpoint: %v", endpoint)
			}

			ip := fmt.Sprintf("%s", body)
			if net.ParseIP(ip) == nil {
				return "", fmt.Errorf("incorrect format return for IPv4 from endpoint: %v", endpoint)
			}

			return ip, nil
		}
	}

	return "", fmt.Errorf("unable to get IP address, ensure machine is connected to the internet")
}
