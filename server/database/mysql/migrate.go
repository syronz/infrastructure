package mysql

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/syronz/infrastructure/server/utils/debug"
	logrus "github.com/Sirupsen/logrus"
)
// Test function for implement select
//
// TODO: Delete this function after implemented inside controllers
func (p *DBT) TestSelect() {

	_, err := DB.Exec("DROP TABLE IF EXISTS users;")
	if err != nil {
		fmt.Println(",,,,,,,,,,,,,,,,,,,,,,,,",err.Error())
	}

	stmt, err := DB.Prepare(`CREATE TABLE users2 (
		id int(6) NOT NULL AUTO_INCREMENT,
		PRIMARY KEY(id)
	) ENGINE=InnoDB DEFAULT CHARSET=latin1;`)
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Table created successfully..")
	}

	rows, err := DB.Query("SELECT * FROM test LIMIT 10")
	if err != nil {
		panic(err.Error())
	}

	type Tag struct {
		ID		int		`json:"id"`
		Name	string	`json:"name"`
	}

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err = rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			panic(err.Error())
		}

		tags = append(tags, tag)

	}

	spew.Dump(">><<<<<<<<<<<<<<",tags)

}


// Create Tables if not exists
func (p *DBT) Migrate(reset bool) {
	logger := logrus.New()
	var err error
	if reset {
		// Drop cities
		_, err = DB.Exec("DROP TABLE IF EXISTS cities;")
		if err != nil {
			debug.Log(err.Error())
		} else {
			logger.Info("Table cities deleted successfully")
		}

		// Drop users
		_, err = DB.Exec("DROP TABLE IF EXISTS users;")
		if err != nil {
			debug.Log(err.Error())
		} else {
			logger.Info("Table users deleted successfully")
		}

	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS cities (
		id int(11) NOT NULL AUTO_INCREMENT,
		governorate varchar(200) NOT NULL,
		city varchar(200) NOT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY city_UNIQUE (city)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;`)
	if err != nil {
		debug.Log(err.Error())
	} else {
		logger.Info("Table cities created successfully")
	}


	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS directors (
		id int(11) NOT NULL AUTO_INCREMENT,
		director varchar(200) NOT NULL,
		detail varchar(200) NOT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY director_UNIQUE (director)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;`)
	if err != nil {
		debug.Log(err.Error())
	} else {
		logger.Info("Table directors created successfully")
	}


	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id int(11) NOT NULL AUTO_INCREMENT,
		name varchar(200) NOT NULL,
		username varchar(200) NOT NULL,
		password varchar(200) NOT NULL,
		role varchar(200) NOT NULL,
		city varchar(500) NOT NULL,
		director varchar(500) NOT NULL,
		language varchar(2) DEFAULT 'en',
		PRIMARY KEY (id),
		UNIQUE KEY username_UNIQUE (username)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;`)
	if err != nil {
		debug.Log(err.Error())
	} else {
		logger.Info("Table users created successfully")
	}

	_, err = DB.Exec(`INSERT IGNORE users(name, username, 
		password, role, city, director, language) VALUES('PLEASE DELETE',
		'a', '$2a$10$rIHWURia91JvlGmgcopll.s7/JS8y0BDB8OLsV8J7ClX2Qi/FEeqG',
		'admin', '*', '*', 'en');`)
	if err != nil {
		debug.Log(err.Error())
	} else {
		logger.Info("Test user created successfully , please delete it later\n Username: a\n Password: a")
	}


	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS activity (
		id int(11) NOT NULL AUTO_INCREMENT,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		event varchar(100) NOT NULL,
		user_id int(11),
		ip varchar(50),
		description varchar(200),
		primary key(id)
	) ENGINE=ARCHIVE;`)
	if err != nil {
		debug.Log(err.Error())
	} else {
		logger.Info("Table activity created successfully")
	}


}
