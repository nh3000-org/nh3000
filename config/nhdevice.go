package config

type DeviceStore struct {
	DSiduuid    string // message id
	DSalias     string // alias
	DShostname  string // hostname
	DSipadrs    string // ip address
	DSmacid     string // macids
	DSnodeuuid  string // unique id
	DSdate      string // message date
	DSsubject   string // message subject
	DSos        string // device os
	DSsequence  uint64 // consumer sequence for secure delete
	DSelementid int    // order in array
}