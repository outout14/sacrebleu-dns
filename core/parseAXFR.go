package core

import (
	"fmt"

	"github.com/miekg/dns"
	"github.com/outout14/sacrebleu-api/api/types"
	"github.com/outout14/sacrebleu-dns/utils"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func parseAXFR(m *dns.Msg) {
	log.Infof("DNS : AXFR query for %s\n", m.Question[0].Name) //Log

	records := utils.GetAllRecords(types.Domain{Fqdn: m.Question[0].Name}) //Get the record in the SQL or Redis database

	for _, record := range records {
		var err error
		var rr dns.RR
		rr, err = dns.NewRR(fmt.Sprintf("%s %v %s %s", record.Fqdn, record.TTL, dns.TypeToString[uint16(record.Type)], record.Content)) //Create the response
		if err == nil {                                                                                                                 //If no err
			m.Answer = append(m.Answer, rr) //Append the record to the response
		} else {
			logrus.Debug(err)
		}
	}
	logrus.Debug(m.Answer)
}
