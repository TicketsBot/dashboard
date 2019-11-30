package database

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	Database gorm.DB
)

func ConnectToDatabase() {
	uri := fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Conf.MariaDB.Username,
		config.Conf.MariaDB.Password,
		config.Conf.MariaDB.Host,
		config.Conf.MariaDB.Database,
	)

	db, err := gorm.Open("mysql", uri)
	if err != nil {
		panic(err)
	}

	db.DB().SetMaxOpenConns(config.Conf.MariaDB.Threads)
	db.DB().SetMaxIdleConns(0)

	db.Set("gorm:table_options", "charset=utf8mb4")
	db.BlockGlobalUpdate(true)

	Database = *db
}
