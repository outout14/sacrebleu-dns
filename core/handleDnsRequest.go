package core

import "github.com/miekg/dns"

//Handle the DNS request using miekg/dns
//Requires dns.ReponseWriter and dns.Msg args
func HandleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {

	//dns.Msg object
	//Will be passed to the parseQuery() function
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	if r.Opcode == dns.OpcodeQuery { //Only respond to dns queries
		parseQuery(m)
	}

	w.WriteMsg(m) //Write the DNS response
}
