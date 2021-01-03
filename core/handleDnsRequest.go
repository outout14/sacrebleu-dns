package core

import (
	"net"

	"github.com/miekg/dns"
	"github.com/outout14/sacrebleu-dns/utils"
)

//HandleDNSRequest : Handle the DNS request using miekg/dns
//Requires dns.ReponseWriter and dns.Msg args
func HandleDNSRequest(w dns.ResponseWriter, r *dns.Msg, conf *utils.Conf) {
	//dns.Msg object
	//Will be passed to the parseQuery() function
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	ip, _, _ := net.SplitHostPort(w.RemoteAddr().String())

	if r.Question[0].Qtype == dns.TypeAXFR {
		if utils.XfrAllowed(ip, conf) {
			parseAXFR(m)
		} else {
			m := new(dns.Msg)
			m.SetRcode(r, dns.RcodeRefused)
			w.WriteMsg(m)
		}

	} else if r.Opcode == dns.OpcodeQuery { //Only respond to dns queries
		parseQuery(m)
	}

	w.WriteMsg(m) //Write the DNS response

}
