package core

import "github.com/miekg/dns"

//Handle the DNS request
func HandleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {

	//dns.Msg object
	//Will be passed to the parseQuery() function
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = true //Less CPU usage (?)

	switch r.Opcode { //Only respond to dns queries
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m) //Write the DNS response
}
