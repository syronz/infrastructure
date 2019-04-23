package mysql

import (
	"github.com/syronz/infrastructure/server/utils/debug"
	logrus "github.com/Sirupsen/logrus"
)
// Test function for implement select
//
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

		// Drop directors
		_, err = DB.Exec("DROP TABLE IF EXISTS directors;")
		if err != nil {
			debug.Log(err.Error())
		} else {
			logger.Info("Table directors deleted successfully")
		}

		// Drop users
		_, err = DB.Exec("DROP TABLE IF EXISTS users;")
		if err != nil {
			debug.Log(err.Error())
		} else {
			logger.Info("Table users deleted successfully")
		}

		// Drop activity
		_, err = DB.Exec("DROP TABLE IF EXISTS activity;")
		if err != nil {
			debug.Log(err.Error())
		} else {
			logger.Info("Table activity deleted successfully")
		}

		// Drop customers
		_, err = DB.Exec("DROP TABLE IF EXISTS customers;")
		if err != nil {
			debug.Log(err.Error())
		} else {
			logger.Info("Table customers deleted successfully")
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

	// Insert data to users
	_, err = DB.Exec(`INSERT IGNORE users(name, username,
		password, role, city, director, language) VALUES('PLEASE DELETE',
		'a', '$2a$10$rIHWURia91JvlGmgcopll.s7/JS8y0BDB8OLsV8J7ClX2Qi/FEeqG',
		'admin', '*', '*', 'en');`)
	if err != nil {
		debug.Log(err.Error())
	} else {
		logger.Info("Test user created successfully , please delete it later\n Username: a\n Password: a")
	}


	// Create activity
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS activity (
		id int(11) NOT NULL AUTO_INCREMENT,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		event varchar(100) NOT NULL,
		user_id int(11),
		ip varchar(50),
		description varchar(200),
		primary key(id)
	) ENGINE=ARCHIVE AUTO_INCREMENT=1000000;`)
	if err != nil {
		debug.Log(err.Error())
	} else {
		logger.Info("Table activity created successfully")
	}


	// Create custoemrs
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS customers (
		id int(11) NOT NULL AUTO_INCREMENT,
		title varchar(20) ,
		name varchar(100) NOT NULL,
		phone1 varchar(20) NOT NULL,
		phone2 varchar(20) ,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		detail varchar(200),
		primary key(id),
		UNIQUE KEY username_UNIQUE (phone1)
	) ENGINE=InnoDB AUTO_INCREMENT=10000;`)
	if err != nil {
		debug.Log(err.Error())
	} else {
		logger.Info("Table customers created successfully")
	}


}


