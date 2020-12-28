package core

import (
	"fmt"

	"github.com/miekg/dns"
	"github.com/outout14/sacrebleu-api/api/types"
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
//Requires dns.ReponseWriter args
func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {

		log.Infof("DNS : Query for %s (type : %v)\n", q.Name, q.Qtype) //Log

		records := utils.GetRecord(types.Record{Fqdn: q.Name, Qtype: q.Qtype}) //Get the record in the SQL or Redis database

		for _, record := range records {
			if record.Content != "" { //If the record is not empty
				log.Infof("DNS : Record found for '%s' => '%s'\n", q.Name, record.Content)                                      //Log the content as INFO
				rr, err := dns.NewRR(fmt.Sprintf("%s %v %s %s", q.Name, record.TTL, dns.TypeToString[q.Qtype], record.Content)) //Create the response
				if err == nil {                                                                                                 //If no err
					m.Answer = append(m.Answer, rr) //Append the record to the response
				}
			} else { //If the record is empty log it as DEBUG
				logrus.Debugf("DNS : No record for '%s' (type '%v')\n", record.Fqdn, record.Qtype)
			}
		}
	}
}
