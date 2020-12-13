package utils

//Structs for configuration
type App struct {
	Port    int
	Ip      string
	Logdir  string
	Logfile bool
}

type Database struct {
	Ip       string
	Port     string
	Username string
	Password string
	Db       string
}

type Redis struct {
	Ip       string
	Port     int
	Password string
	Db       int
	Ttl      int
}

type Conf struct {
	App_mode string
	App
	Database
	Redis
}

type Domain struct {
	ID           int `json:"id"`
	FriendlyName string
	Fqdn         string
	OwnerId      int
	LastEdit     string
}

type Record struct {
	Id       int
	DomainId int
	Fqdn     string
	Content  string
	Type     int
	Qtype    uint16
	TTL      int
}
