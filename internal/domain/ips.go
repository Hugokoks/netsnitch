package domain

import (
	"fmt"
	"net"
)


func ParseCIDR(cidr string) ([]net.IP,error){


	ip,ipNet,err := net.ParseCIDR(cidr)

	if err != nil{

		return nil, fmt.Errorf("Invalid CIDR: %w",err)
	}

	var ips []net.IP

	////first IP in network
	firstNetworkIP := ip.Mask(ipNet.Mask)
	
	for currIP := firstNetworkIP; ipNet.Contains(currIP); incIP(currIP){
		// Copy the IP because net.IP is a mutable byte slice.
		// currIP is reused and modified in each iteration.
		ipCopy := make(net.IP,len(currIP))
		copy(ipCopy,currIP)
		ips = append(ips,ipCopy)

	}
	return ips, nil
}

/////increment IP by 1
func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] != 0 {
			break
		}
	}
}
