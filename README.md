# LightHouse

LightHouse is a DNS authoritative nameserver made in Go with ``github.com/miekg/dns`` library. 

Records are stored in a MySQL Database and cached using REDIS. 

This software is currently in development and NOT ready for production.

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