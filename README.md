# Sacrebleu
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Foutout14%2Fsacrebleu-dns.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Foutout14%2Fsacrebleu-dns?ref=badge_shield)


Sacrebleu is a DNS authoritative nameserver made in Go with ``github.com/miekg/dns`` library. 

Records are stored in a SQL Database and cached using Redis.

This software requires a SQL server and a Redis server.

This software is currently in development and NOT ready for production.

## Installation 

### Automated install
Each release contains a RPM and DEB package. 
- [Download the latest package](https://github.com/outout14/sacrebleu-dns/releases/latest) compatible with your system 
- Edit the config file in ``/etc/sacrebleu/config-dns.ini`` 
- Launch the database migration with ``sacrebleu-dns -sqlmigrate -config /etc/sacrebleu/config-dns.ini``
- Start the service ``systemctl start sacrebleu-dns`` 

### Manual Install
- [Download the latest release](https://github.com/outout14/sacrebleu-dns/releases/latest) compatible with your system 
- [Download](https://raw.githubusercontent.com/outout14/sacrebleu-dns/main/extra/config.ini.example) and edit the configuration file  
- Launch the database migration with ``sacrebleu-dns -sqlmigrate -config [config path].ini``
- [Download](https://raw.githubusercontent.com/outout14/sacrebleu-dns/main/extra/sacrebleu-dns.service) and edit the systemd service and put it in ``/lib/systemd/system/`` 
- Start the service ``systemctl start sacrebleu-dns`` 

## Arguments 
You can show theses informations using ``./sacrebleu-dns -h``.
``` 
   -config string
        the patch to the config file (default "extra/config.ini.example")
  -sqlmigrate
        initialize / migrate the database
``` 

## Configuration 
Variables names are case sensitives.
|Variable name|Type|Example|Informations|
|--|--|--|--|
| AppMode | string|``"production"``|Anything different than ``production`` will show debug messages
| App | Section |
|IP|string|``":"``|IP address on which the DNS server must listen. Blank to listen on all IPs 
|Port|int|``5353``|Port on which the DNS server must listen
|Logfile|bool|``true``|Enable or disable file logs.
|Database|Section|
|IP|string|``"127.0.0.1"``|SQL Database IP 
|Username|string|``"sacrebleu"``|SQL Database Username
|Password|string|``"superSecretPassword"``|SQL Database Password
|Port|string|``"5432"``|SQL Database port (``"5432"`` for postgres or ``"3306"`` for MySQL by default)
|Db|string|``"sacrebleudatabase"``|SQL Database Name 
|Type|string|``"postgresql"``|SQL Database type. ``"postgresql"`` or ``"mysql"`` (anything different than ``"postgresql"`` will rollback to ``"mysql"``
|Redis|Section
|IP|string|``"127.0.0.1"``|Redis Database IP 
|Password|string|``""``|Redis Database Password
|Port|int|``6379``|Redis Database port
|Db|int|``0``|Redis Database ID
|Ttl|int|``10``|Redis Time To Live (in seconds)

## What is working 
- Read records (stricts & wildcard) from MySQL
- Read and write records (stricts & wildcard) in Redis
- Recursive wildcard for reverse DNS (IPv6 only) 
- Generate dynamic reverse DNS (IPv6 only)
- Respond to all requested MySQL queries 

## ToDo
- Recursive wildcard for reverse DNS (IPv4 part) 
- Generate dynamic reverse DNS (IPv4 part)
- XFR 
- DNSSEC 
- Unit tests 


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Foutout14%2Fsacrebleu-dns.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Foutout14%2Fsacrebleu-dns?ref=badge_large)
