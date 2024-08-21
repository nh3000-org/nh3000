package config

type CommandStore struct {
	CMDalias    string // alias
	CMDtype     int    // 1 = bash, 2 = bat. 3 = snmp
	CMDinterval string // how often
	MScommand   string // command to execute
	MSexpected  string // expected result from command
}
