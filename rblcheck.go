package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
)

import(
	rblnet "./lib/net"
	"./lib/rbl"
)

var rbls = []string{
	"zen.spamhaus.org.",
	"bl.spamcop.net.",
	"short.rbl.jp.",
	"dnsbl.sorbs.net.",
	"cbl.abuseat.org.",
	"abuse.rfc-ignorant.org.",
	"b.barracudacentral.org.",
	"db.wpbl.info.",
	"black.junkemailfilter.com.",
	"bl.mailspike.net.",
	"psbl.surriel.com.",
	"ubl.unsubscore.com.",
}

var verbose verboseT = false

type verboseT bool

func (v verboseT) Printf(format string, args ...interface{}) {
	if v {
		fmt.Printf(format, args...)
	}
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s addr...\n\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	usage := "verbose"

	v := (*bool)(&verbose)
	flag.BoolVar(v, "verbose", false, usage)
	flag.BoolVar(v, "v", false, usage + " (shorthand)")
}

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	ips := make([]net.IP, 0, len(flag.Args()))
	for _, v := range flag.Args() {
		verbose.Printf("lookup ip %v\n", v)
		addrs, err := rblnet.LookupIP(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid address %s\n", v)
			os.Exit(1)
		}
		ips = append(ips, addrs...)
	}

	verbose.Printf("check %d addresses\n", len(ips))

	stdout := make(chan string)
	stderr := make(chan string)
	done := make(chan bool)

	for _, v := range rbls {
		rbl := rbl.NewRBL(v)
		check(rbl, ips, (chan<- string)(stdout), (chan<- string)(stderr), (chan<- bool)(done))
	}

	remain := len(rbls)
	exit := 0

	loop:
	for {
		select {
		case str := <-stdout:
			fmt.Fprint(os.Stdout, str)
		case str := <-stderr:
			fmt.Fprint(os.Stderr, str)
		case ret := <-done:
			if !ret { exit = 1 }
			remain--
			if remain == 0 {
				break loop
			}
		}
	}

	os.Exit(exit)
}

func check(rbl *rbl.RBL, ips []net.IP, stdout chan<- string, stderr chan<- string, done chan<- bool) {
	go func() {
		listed := false
		defer func(){ done<- !listed }()
		for _, ip := range ips {
			if verbose {
				stdout<- fmt.Sprintf("checking %s in %s\n", ip, rbl.Zone)
			}
			res, _ := rbl.LookupRBL(ip)
			if res.Listed == true {
				str := fmt.Sprintf("%s listed in %s", res.Ip, res.Zone)
				if t := res.Text; t != "" {
					t = strings.Replace(t, "\n", " ", -1)
					str += fmt.Sprintf(": %s", t)
				}
				stderr<- fmt.Sprintln(str)
				listed = true
			}
		}
	}()
}
