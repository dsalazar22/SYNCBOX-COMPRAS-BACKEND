package libs

type InfoUser struct {
	IdUser   string `json:"iduser"`
	Username string `json:"user"`
	Email    string `json:"email"`
}

type ConfigConnection struct {
	Addr           string
	Dest           string
	Vars           string
	Pass           string
	Db             int
	DbTS           int
	Pg             string
	UserURL        string
	ModulesActives string
	Files          string
	HourStart      string
}

var StrConn ConfigConnection
