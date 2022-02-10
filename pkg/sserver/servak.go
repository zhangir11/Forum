package sserver

import (
	"forum/pkg/databasesql"
)

//Servak ...
var Servak *Server

//Db ...
var Db *databasesql.DatabaseSQL

//DeclareServak ...
func DeclareServak(s *Server) {
	Servak = s
}

//DeclareDB ...
func DeclareDB(s *databasesql.DatabaseSQL) {
	Db = s
}
