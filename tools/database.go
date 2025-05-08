package tools

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/responsible-api/responsible-auth/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const fmtMysqlDBString = "%s:%s@tcp(%s:%d)/%s"

func NewDatabase() (*gorm.DB, error) {
	db, err := DBCon()
	if err != nil {
		log.Fatalf("DB connection start failure error")
		return nil, err
	}
	return db, nil
}

func DBCon() (*gorm.DB, error) {
	c := config.ConfigDB()
	l := logger.Default.LogMode(logger.Silent)
	if c.Debug {
		l = logger.Default.LogMode(logger.Info)
	}

	dsn := fmt.Sprintf(fmtMysqlDBString, c.Username, c.Password, c.Host, c.Port, c.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: l,
	})
	if err != nil {
		log.Fatalf("DB connection start failure error")
		return nil, err
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("DB connection start failure error")
			return
		}
		sqlDB.Close()
	}()
	return db, nil
}
