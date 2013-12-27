package goat

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Connect to MySQL database
func DbConnect() (*sqlx.DB, error) {
	return sqlx.Connect("mysql", fmt.Sprintf("%s:%s@/%s", "goat", "goat", "goat"))
}

func DbManager(dbDoneChan chan bool) {
	// Storage handler instances
	mapDb := new(MapDb)
	sqlDb := new(SqlDb)

	// channels
	sqlRequestChan := make(chan Request)
	mapRequestChan := make(chan Request, 100)

	// Shutdown function
	go func(dbDoneChan chan bool, mapDb *MapDb, sqlDb *SqlDb) {
		// Wait for shutdown
		Static.ShutdownChan <- <-Static.ShutdownChan
		Static.ShutdownChan <- true

		if Static.Config.Map {
			mapDb.Shutdown()
		}
		if Static.Config.Sql {
			sqlDb.Shutdown()
		}

		dbDoneChan <- true
	}(dbDoneChan, mapDb, sqlDb)

	if Static.Config.Map && Static.Config.Sql {
		for {
			select {
			case hold := <-Static.RequestChan:
				if hold.Data == nil {
					mapRequestChan <- hold
				} else {
					mapRequestChan <- hold
					sqlRequestChan <- hold
				}
			case hold := <-Static.PersistentChan:
				sqlRequestChan <- hold
			}
		}
	} else if Static.Config.Map {
		for {
			select {
			case hold := <-Static.RequestChan:
				mapRequestChan <- hold
			}
		}
	} else if Static.Config.Sql {
		for {
			select {
			case hold := <-Static.RequestChan:
				sqlRequestChan <- hold
			}
		}
	} else {
		Static.LogChan <- "No database in use."
	}
}

// DbHandler interface method HandleDb defines a database handler which handles requests
type DbHandler interface {
	Read(chan Request)
	Write(chan Request)
	Shutdown()
}
