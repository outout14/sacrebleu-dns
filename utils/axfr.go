package utils

import (
	"github.com/outout14/sacrebleu-api/api/types"
)

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
	results = append([]types.Record{soa}, results...)

	return results
}
