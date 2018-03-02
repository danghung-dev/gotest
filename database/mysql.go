package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"gotest/config"
)

type MySQLDB struct {
	*gorm.DB
}

func NewMySQLDB(dbCfg config.MySQLConfig) (*MySQLDB, error) {
	//DSN := fmt.Sprintf("%s:%s@unix(/tmp/mysql.sock)/%s?parseTime=true", dbCfg.Username, dbCfg.Password, dbCfg.DatabaseName)
	dataSourceName := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true", dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DatabaseName, dbCfg.Encoding)
	db, err := gorm.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &MySQLDB{db}, nil
}