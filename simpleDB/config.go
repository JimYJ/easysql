package simpleDB

import (
	"sync"
)

var (
	mysqldb *MysqlDB
	once    sync.Once
)

/* mysql db param */
const (
	MysqlDBHost = "0.0.0.0"
	MysqlDBuser = "root"
	MysqlDBpass = ""
	MysqlDBPort = 3306
	MysqlDBName = ""
)

func GetMysqlConn() (*MysqlDB, error) {
	var err error
	once.Do(func() {
		mysqldb = &MysqlDB{MysqlDBHost, MysqlDBuser, MysqlDBName, MysqlDBpass, MysqlDBPort, nil, nil, nil}
		err = mysqldb.Conn(75, 75)
	})
	return mysqldb, err
}
