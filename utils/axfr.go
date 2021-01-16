package utils

import (
	"github.com/outout14/sacrebleu-api/api/types"
)

//XfrAllowed : check if the IP is allowed to perform XFR requests
func XfrAllowed(remoteIP string, conf *Conf) bool {
	for _, ip := range conf.DNS.XfrIPs {
		if ip == "*" {
			return true
		}
		if ip == remoteIP {
			return true
		}
	}
	return false
}

//GetAllRecords : Retrive all records for a domain
func GetAllRecords(d types.Domain) []types.Record {
	results := []types.Record{}

	//Get domain id from his fqdn
	err := d.GetDomainByFqdn(db)
	if err != nil {
		return results
	}

	results, err = d.GetDomainRecords(db, -1, -1)
	DbgErr(err)

	soa, err := d.GetSOA(db)
	DbgErr(err)

	results = append([]types.Record{soa}, results...)

	return results
}
