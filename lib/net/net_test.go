package util

import(
	"github.com/stretchr/testify/assert"
	"testing"
	"net"
)

func Test_ipToLong(t *testing.T) {
	ip := net.ParseIP("192.168.0.1")
	l, _ := ipToLong(ip)
	assert.IsType(t, uint32(0), l, "should be a type of uint32")
}

func Test_longToIP(t *testing.T) {
	var l uint32

	l = uint32(1)
	ip := longToIP(l)
	assert.Equal(t, "0.0.0.1", ip.String(), "should be '0.0.0.1'")

	l += 255
	ip = longToIP(l)
	assert.Equal(t, "0.0.1.0", ip.String(), "should be '0.0.1.0'")
}

func Test_HostAddr(t *testing.T) {
	var ipaddr string
	var ipnet *net.IPNet
	var addrs []net.IP

	ipaddr = "192.168.0.0/24"
	_, ipnet, _ = net.ParseCIDR(ipaddr)
	addrs, _ = HostAddr(ipnet)
	assert.Equal(t, 254, len(addrs), "should be equal 254") /* excluded network-address and broadcast-address */
	assert.Equal(t, "192.168.0.1", addrs[0].String(), "should be '192.168.0.1'")
	assert.Equal(t, "192.168.0.254", addrs[len(addrs) - 1].String(), "should be '192.168.0.254'")

	ipaddr = "192.168.0.1/32"
	_, ipnet, _ = net.ParseCIDR(ipaddr)
	addrs, _ = HostAddr(ipnet)
	assert.Equal(t, 1, len(addrs), "should be equal 1")
	assert.Equal(t, "192.168.0.1", addrs[0].String(), "should be '192.168.0.1'")

	ipaddr = "2001:DB8::/48"
	_, ipnet, _ = net.ParseCIDR(ipaddr)
	_, err := HostAddr(ipnet)
	assert.Error(t, err, "occurred with IPv6 network")
}

func Test_LookupIP(t *testing.T) {
	var T []net.IP
	var host string
	var ip []net.IP

	host = "192.168.0.1"
	ip, _ = LookupIP(host)
	assert.IsType(t, T, ip, "should be a type of []net.IP")
	assert.Equal(t, 1, len(ip), "should be equal 1")
	assert.Equal(t, "192.168.0.1", ip[0].String(), "should be equal '192.168.0.1'")

	host = "192.168.0.102/32"
	ip, _ = LookupIP(host)
	assert.IsType(t, T, ip, "should be a type of []net.IP")
	assert.Equal(t, 1, len(ip), "should be equal 1")
	assert.Equal(t, "192.168.0.102", ip[0].String(), "should be equal '192.168.0.102'")

	host = "192.168.0.0/24"
	ip, _ = LookupIP(host)
	assert.IsType(t, T, ip, "should be a type of []net.IP")
	assert.Equal(t, 254, len(ip), "should be equal 1")
	assert.Equal(t, "192.168.0.1", ip[0].String(), "should be '192.168.0.1'")
	assert.Equal(t, "192.168.0.254", ip[len(ip) - 1].String(), "should be '192.168.0.254'")

	host = "google.com"
	ip, _ = LookupIP(host)
	assert.IsType(t, T, ip, "should be a type of []net.IP")
	assert.True(t, len(ip) > 0, "should be true")

	host = "invalid-host(address|name)"
	ip, err := LookupIP(host)
	assert.IsType(t, T, ip, "should be a type of []net.IP")
	assert.True(t, len(ip) == 0, "should be true")
	assert.Error(t, err, "occurred with invalid host")
}
