package utils

import (
	"github.com/m9rco/exile/kernel/common"
	"net"
)

func GetLocalIP() (ipv4 string, err error) {
	var (
		address []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet
		isIpNet bool
	)
	// get add address
	if address, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// get first
	for _, addr = range address {
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// continue IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String() // 192.168.1.1
				return
			}
		}
	}
	err = common.ERROR_NO_LOCAL_IP_FOUND
	return
}
