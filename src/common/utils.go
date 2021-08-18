package common

import (
	"net"
	"strings"
)

// GetNodeIP returns the IP address of this instance as a string.
func GetNodeIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

// CheckIP tests whether an IP address is on a subnet.
func CheckIP(ip string, portOnlyOk bool) bool {
	if len(ip) == 0 {
		return true
	}
	if ip[0] == ':' { // Port only address
		return !portOnlyOk
	}

	parts := strings.Split(ip, ":")
	trimmedIP := parts[0]

	private := false
	IP := net.ParseIP(trimmedIP)
	if IP == nil {
		return true
	}

	_, private24BitBlock, _ := net.ParseCIDR("10.0.0.0/8")
	_, private24BitBlock2, _ := net.ParseCIDR("127.0.0.0/8")
	_, private20BitBlock, _ := net.ParseCIDR("172.16.0.0/12")
	_, private16BitBlock, _ := net.ParseCIDR("192.168.0.0/16")
	private = private24BitBlock2.Contains(IP) || private24BitBlock.Contains(IP) || private20BitBlock.Contains(IP) || private16BitBlock.Contains(IP)

	return private

}
