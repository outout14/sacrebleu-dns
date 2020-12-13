package core

import (
	"fmt"

	"github.com/miekg/dns"
	"github.com/outout14/sacrebleu-dns/utils"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

/*
	Qtype memo
	A = 1
	NS = 2
	PTR = 12
	TXT = 16
	AAAA = 28
*/

//Function called by handleDnsRequest to parse the query from records
func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {

		log.Infof("DNS : Query for %s (type : %v)\n", q.Name, q.Qtype) //Log

		record := utils.GetRecord(utils.Record{Fqdn: q.Name, Qtype: q.Qtype})

		if record.Content != "" { //If the record is found, return it
			log.Infof("DNS : Record found for '%s' => '%s'\n", q.Name, record.Content)
			rr, err := dns.NewRR(fmt.Sprintf("%s %v %s %s", q.Name, record.TTL, dns.TypeToString[q.Qtype], record.Content)) //Create the response
			if err == nil {                                                                                                 //If no err
				m.Answer = append(m.Answer, rr)
			}
		} else {
			logrus.Debugf("DNS : No record for '%s' (type '%v')\n", record.Fqdn, record.Qtype)
		}

	}
}
