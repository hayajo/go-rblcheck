package main

import(
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
	"os"
	"net"
)

func Test_NewRBL(t *testing.T) {
	var rbl *RBL

	rbl = NewRBL("rbl.example.com")
	assert.IsType(t, new(RBL), rbl, "should be a type of *RBL")
}

func Test_auery(t *testing.T) {
	var zone string = "rbl.example.com"
	var rbl *RBL = NewRBL(zone)
	var ip net.IP
	var query string

	ip = net.ParseIP("192.168.0.1")
	query, _ = rbl.query(ip)
	assert.Equal(t, "1.0.168.192" + "." + zone, query, zone, "should be '1.0.168.192.rbl.example.com'")

	ip = net.ParseIP("2001:DB8::")
	query, err := rbl.query(ip)
	assert.Error(t, err, "occurred with IPv6 address")
}

func Test_LookupRBL(t *testing.T) {
	var zone string = "zen.spamhaus.org" /* TODO mocking */
	var rbl *RBL = NewRBL(zone)
	var ip net.IP
	var res RBLResult
	var err error

	ip = net.ParseIP("127.0.0.1")
	res, err = rbl.LookupRBL(ip)
	if assert.False(t, res.Listed, "occurred with unlisted host") != true {
		fmt.Fprintln(os.Stderr, res, err)
	}

	ip = net.ParseIP("127.0.0.2")
	res, err = rbl.LookupRBL(ip)
	if assert.True(t, res.Listed, "should be true") != true {
		fmt.Fprintln(os.Stderr, res, err)
	}
	assert.NotEmpty(t, res.Text, "should be not empty")
}
