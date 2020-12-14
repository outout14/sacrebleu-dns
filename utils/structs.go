package utils

//Struct for App (dns server) configuration in the config.ini file
type App struct {
	Port    int
	Ip      string
	Logdir  string
	Logfile bool
}

//Struct for SQL Database configuration in the config.ini file
type Database struct {
	Ip       string
	Port     string
	Username string
	Password string
	Db       string
}

//Struct for Redis  Database configuration in the config.ini file
type Redis struct {
	Ip       string
	Port     int
	Password string
	Db       int
	Ttl      int
}

//Struct for the whole config.ini file when it will be parsed by go-ini
type Conf struct {
	App_mode string
	App
	Database
	Redis
}

//Struct for a Domain (not used currently).
type Domain struct {
	ID           int `json:"id"`
	FriendlyName string
	Fqdn         string
	OwnerId      int
	LastEdit     string
}

//Struct for a domain record
//Defined by it's ID, DomainID (parent domain), Fqdn (or name), Content (value of the record), Type (as Qtype/int), TTL (used only for the DNS response and not the Redis TTL)
type Record struct {
	Id       int
	DomainId int
	Fqdn     string
	Content  string
	Type     int
	Qtype    uint16
	TTL      int
}
