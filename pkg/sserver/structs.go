package sserver

import (
	"forum/pkg/funcs"
)

// Server ...
type Server struct {
	NbrOfConnections int
	WebPort          string `json:"port"`
	Database         funcs.Database
	NameOfDatabase   string `json:"namedb"`
	TypeOfDatabase   string `json:"typedb"`
	MaxConnections   int    `json:"maxconnection"`
	MaxTime          int    `json:"maxtime"`
}

//NewServer ...
func NewServer() *Server {
	return &Server{}
}
