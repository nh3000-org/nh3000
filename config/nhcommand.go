package config

type CommandStore struct {
	CMDalias    string // alias
	CMDtype     int    // 1 = bash, 2 = bat. 3 = snmp
	CMDinterval string // how often
	CMDcommand  string // command to execute
	CMDexpected string // expected result from command
	CMDtimeout  int    // time to wait in seconds, 0 one ti,e

}
