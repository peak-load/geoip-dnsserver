/// Copyright (c) 2021, peak-load
/// All rights reserved.

/// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package main

import (
	"log"
	"net"
	"regexp"
	"strconv"

	"github.com/miekg/dns"
	"github.com/oschwald/geoip2-golang"
)

type Records struct {
	Name                         string
	Address                      string
	AutonomousSystemNumber       uint
	AutonomousSystemOrganization string
	ISP                          string
	Organization                 string
}

func parseQuery(m *dns.Msg, addressOfRequester string) {
	r, _ := regexp.Compile("^geoip.")
	for _, q := range m.Question {
		if r.MatchString(q.Name) {
			switch q.Qtype {
			case dns.TypeTXT:
				log.Printf("Query for %s\n", addressOfRequester)
				geoip := net.ParseIP(addressOfRequester)

				db, err := geoip2.Open("GeoIP2-ISP.mmdb")
				if err != nil {
					log.Fatal(err)
				}
				defer db.Close()
				record, err := db.ISP(geoip)
				if err != nil {
					log.Fatal(err)
				}
				var txt string = "ip='" + geoip.String() + "',asn='" + strconv.FormatUint(uint64(record.AutonomousSystemNumber), 10) + "',asn_organization='" + record.AutonomousSystemOrganization + "',isp='" + record.ISP + "',organization='" + record.Organization + "'"
				if record.ISP == "" {
					txt = "Sorry, but your IP is not known."
				}

				rr := &dns.TXT{Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 3600}, Txt: []string{txt}}
				if err == nil {
					m.Answer = append(m.Answer, rr)
				} else {
					log.Fatalf("Failed to conftrct TXT record %s\n ", err.Error())
				}
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	addressOfRequester, _, err := net.SplitHostPort(w.RemoteAddr().String())
	if err != nil {
		log.Fatalf("Failed to get hostname %s\n ", err.Error())
	}
	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m, addressOfRequester)
	}
	w.WriteMsg(m)
}

func main() {
	dns.HandleFunc(".", handleDnsRequest)
	// specify port number
	port := 53
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
