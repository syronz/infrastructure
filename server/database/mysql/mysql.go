// MySQL package.
//
// Just struct and connect function, extra queries will be executed inside controllers
package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/syronz/infrastructure/server/app"
	"github.com/syronz/infrastructure/server/utils/debug"
)

type DBT struct{}

var DBTi DBT

var DB *sql.DB

// Connect to the database with information that comes from app.toml
func (p *DBT) Connect() {
	var err error
	DB, err = sql.Open("mysql", app.Config.Mysql.User + ":" +
	app.Config.Mysql.Password + "@tcp(" + app.Config.Mysql.Host +
	":"+ app.Config.Mysql.Port +")/" + app.Config.Mysql.Database)
	if err != nil {
		debug.Log(err)
		panic(err.Error())
	}
}


