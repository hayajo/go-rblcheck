package main

import(
	"strconv"
	"strings"
	"fmt"
	"net"
	"regexp"
)

type RBL struct {
	Zone string
}

type RBLResult struct {
	Listed bool
	Ip net.IP
	Zone string
	Text string
}

func NewRBL(zone string) *RBL {
	return &RBL{Zone: zone}
}

func (rbl *RBL) LookupRBL(ip net.IP) (RBLResult, error) {
	res := RBLResult{
		Listed: false,
		Ip: ip,
		Zone: rbl.Zone,
		Text: "",
	}

	query, err := rbl.query(ip)
	if err != nil {
		return res, err
	}

	addrs, err := net.LookupHost(query)
	reg, _ := regexp.Compile("(?i:no such host)")
	if err != nil && reg.FindStringIndex(err.Error()) == nil {
		return res, err
	}

	if len(addrs) > 0 {
		res.Listed = true
		text, err := net.LookupTXT(query)
		if err == nil {
			res.Text = strings.Join(text, "\n")
		}
	}

	return res, nil
}

func (rbl *RBL) query(ip net.IP) (string, error) {
	v4 := ip.To4()
	if v4 == nil {
		return "", fmt.Errorf("Support only IPv4: %v", ip)
	}

	inverted := make([]string, len(v4))
	for i := 0; i < len(v4); i++ {
		inverted[i] = strconv.Itoa(int(v4[len(v4) - i - 1]))
	}

	query := strings.Join(inverted, ".") + "." + rbl.Zone

	return query, nil
}
