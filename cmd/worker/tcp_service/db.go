package tcp_service

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strconv"
)

var (
	driverName   = "mysql"
	databasePort = 3306
	databaseName = "entrytask"
	dataSource   = "root:todd123456@tcp(localhost:" + strconv.Itoa(databasePort) + ")/" + databaseName
)

func InitDB() *sqlx.DB {

	DBEngine, err := sqlx.Connect(driverName, dataSource)
	if err != nil {
		logrus.Panic(err)
	}
	DBEngine.SetMaxIdleConns(100)
	DBEngine.SetMaxOpenConns(200)
	logrus.Infof("init database %s", databaseName)
	return DBEngine
}
