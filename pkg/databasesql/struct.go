package databasesql

import "database/sql"

// DatabaseSQL ...
type DatabaseSQL struct {
	DatabaseAddress string
	Db              *sql.DB
}

//NewDatabaseSQL ...
func NewDatabaseSQL() *DatabaseSQL {
	return &DatabaseSQL{}
}

// Sessions ...
type Sessions struct {
	UserID       int64
	SessionName  string
	SessionValue string
}
