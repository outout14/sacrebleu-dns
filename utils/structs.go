package utils

//App : Struct for App (dns server) configuration in the config.ini file
type App struct {
	Port    int
	IP      string `ini:"IP"`
	Logdir  string
	Logfile bool
}

//Database : Struct for SQL Database configuration in the config.ini file
type Database struct {
	Host     string `ini:"Host"`
	Port     string
	Username string
	Password string
	Db       string `ini:"DB"`
	Type     string
}

//Redis : Struct for Redis  Database configuration in the config.ini file
type Redis struct {
	IP       string `ini:"IP"`
	Port     int
	Password string
	Db       int `ini:"DB"`
	TTL      int `ini:"TTL"`
}

//DNS : Struct for XFR and NS
type DNS struct {
	XfrIPs      []string
	Nameservers []string
}

//Conf : Struct for the whole config.ini file when it will be parsed by go-ini
type Conf struct {
	AppMode string `ini:"app_mode"`
	App
	Database
	Redis
}
