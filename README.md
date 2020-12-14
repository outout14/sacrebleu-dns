# Sacrebleu
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Foutout14%2Fsacrebleu-dns.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Foutout14%2Fsacrebleu-dns?ref=badge_shield)


Sacrebleu is a DNS authoritative nameserver made in Go with ``github.com/miekg/dns`` library. 

Records are stored in a MySQL Database and cached using REDIS. 

This software requires MySQL and REDIS servers.

This software is currently in development and NOT ready for production.

## Initialization 
- Move ``./extra/config.ini.example`` to ``./config.ini`` and edit it.
- Import to your configured database the SQL structure dump from ``./extra/db.sql`` 
- ``make build``

## What is working 
- Read records (stricts & wildcard) from MySQL
- Read and write records (stricts & wildcard) in REDIS
- Recursive wildcard for reverse DNS (IPv6 only) 
- Generate dynamic reverse DNS (IPv6 only)
- Respond to all requested MySQL queries 

## ToDo 
- Recursive wildcard for reverse DNS (IPv4 part) 
- Generate dynamic reverse DNS (IPv4 part)
- XFR 
- DNSSEC 
- Unit tests 
- CI with auto packaging
- Optimization


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Foutout14%2Fsacrebleu-dns.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Foutout14%2Fsacrebleu-dns?ref=badge_large)