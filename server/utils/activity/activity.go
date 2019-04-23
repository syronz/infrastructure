package activity

import (
	"github.com/syronz/infrastructure/server/database/mysql"
	"github.com/syronz/infrastructure/server/utils/debug"
	"github.com/syronz/infrastructure/server/app"

	"net"
	"github.com/syronz/ozzo-routing"
)

// Log All activity inside application
func Log(c *routing.Context, event string, description string) {
	userId := app.GetRequestScope(c).UserID()
	ip,_,_ := net.SplitHostPort(c.Request.RemoteAddr)
	debug.Log(event, description)

	db := mysql.DB
	stmt, err := db.Prepare(`INSERT INTO activity(event, user_id, ip, description) VALUES(?, ?, ?, ?)`)
	defer stmt.Close()
	if err != nil {
		debug.Log(err.Error())
		return
	}
	_, err = stmt.Exec(event, userId, ip, description)
	if err != nil {
		debug.Log(err.Error())
		return
	}

}


