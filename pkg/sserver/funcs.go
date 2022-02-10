package sserver

import (
	"encoding/json"
	"forum/pkg/databasesql"
	"forum/pkg/funcs"
	"log"
	"net/http"
	"os"
)

//CreateServer ...
func CreateServer() *Server {
	SServer := NewServer()

	if jsn, err := os.ReadFile("configs/config.json"); err == nil {
		if err1 := json.Unmarshal(jsn, SServer); err1 != nil {
			log.Println("Failed on configure: ", err1.Error())
			return nil
		}
	} else {
		log.Println("Failed on configure: ", err.Error())
		return nil
	}
	if SServer.TypeOfDatabase == "sql" {
		var temp *databasesql.DatabaseSQL = databasesql.NewDatabaseSQL()
		var Db funcs.Database = temp
		DeclareDB(temp)
		err1 := Db.Initiate(SServer.NameOfDatabase)
		if err1 != nil {
			log.Println(err1.Error())
		}
		SServer.Database = Db
	}
	return SServer
}

//CheckIfUserAuthorized ...
func CheckIfUserAuthorized(r *http.Request) (*funcs.User, error) {
	c, err := r.Cookie("authenticated")
	if err == nil {
		u, err := Servak.Database.GetUserByCookie(c.Value)
		if err != nil {
			return funcs.NewUser(), err
		}
		return u, nil
	}
	return funcs.NewUser(), err
}
