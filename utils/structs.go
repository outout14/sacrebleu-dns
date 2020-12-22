package utils

//App : Struct for App (dns server) configuration in the config.ini file
type App struct {
	Port    int
	IP      string
	Logdir  string
	Logfile bool
}

//Database : Struct for SQL Database configuration in the config.ini file
type Database struct {
	IP       string
	Port     string
	Username string
	Password string
	Db       string
	Type     string
}

//Redis : Struct for Redis  Database configuration in the config.ini file
type Redis struct {
	IP       string
	Port     int
	Password string
	Db       int
	TTL      int
}

//Conf : Struct for the whole config.ini file when it will be parsed by go-ini
type Conf struct {
	AppMode string
	App
	Database
	Redis
}

//Domain : Struct for a Domain (not used currently).
type Domain struct {
	ID           int `json:"id"`
	FriendlyName string
	Fqdn         string
	OwnerID      int
	LastEdit     string
}

//Record : Struct for a domain record
//Defined by it's ID, DomainID (parent domain), Fqdn (or name), Content (value of the record), Type (as Qtype/int), TTL (used only for the DNS response and not the Redis TTL)
type Record struct {
	ID       int
	DomainID int
	Fqdn     string
	Content  string
	Type     int
	Qtype    uint16
	TTL      int
}
