package main

import(
	"encoding/binary"
	"net"
	"math"
	"fmt"
)

func ipToLong(ip net.IP) (long uint32, err error) {
	if ip.To4() == nil {
		return 0, fmt.Errorf("Support only IPv4 address")
	}
	return binary.BigEndian.Uint32(ip), nil
}

func longToIP(long uint32) net.IP {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(long))
	return net.IP(bytes).To4()
}

func HostAddr(ipn *net.IPNet) (addrs []net.IP, err error) {
	if ipn.IP.To4() == nil {
		return nil, fmt.Errorf("Support only IPv4: %v", ipn)
	}

	masklen, bitlen := ipn.Mask.Size()

	if masklen > 30 {
		return []net.IP{ipn.IP}, nil
	}
	
	nwLong, _ := ipToLong(ipn.IP) /* convert network-address */

	hostlen := int(math.Pow(2, float64(bitlen - masklen)) - 2)
	hostenum := make([]net.IP, hostlen)

	for i := 0; i < hostlen; i++ {
		host := longToIP(nwLong + uint32(i + 1))
		hostenum[i] = host
	}

	return hostenum, nil
}

func LookupIP(host string) (addrs []net.IP, err error) {
	ip := net.ParseIP(host)
	if ip == nil || ip.To4() == nil {

		ip, ipnet, err := net.ParseCIDR(host)
		if err != nil || ip.To4() == nil {

			hosts, err := net.LookupHost(host)
			if err != nil {
				return nil, err
			}

			addrs := make([]net.IP, 0, len(hosts))
			for _, v := range hosts {
				ip = net.ParseIP(v)
				ip = ip.To4()
				if ip != nil {
					addrs = append(addrs, ip)
				}
			}

			return addrs, nil
		}

		addrs, _ := HostAddr(ipnet)
		return addrs, nil
	}

	return []net.IP{ip}, nil
}
