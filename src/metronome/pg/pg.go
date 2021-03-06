package pg

import (
	"sync"

	"github.com/spf13/viper"
	"gopkg.in/pg.v5"
)

type db struct {
	DB *pg.DB
}

var d *db
var onceDB sync.Once

// DB get a database instance
func DB() *pg.DB {
	onceDB.Do(func() {
		database := pg.Connect(&pg.Options{
			Addr:     viper.GetString("pg.addr"),
			User:     viper.GetString("pg.user"),
			Password: viper.GetString("pg.password"),
			Database: viper.GetString("pg.database"),
		})

		// Ensure tables
		if _, err := database.Exec(string(MustAsset("users.sql"))); err != nil {
			panic(err)
		}
		if _, err := database.Exec(string(MustAsset("tasks.sql"))); err != nil {
			panic(err)
		}

		d = &db{
			DB: database,
		}
	})
	return d.DB
}
