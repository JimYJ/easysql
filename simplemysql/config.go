package simplemysql

import (
	"sync"
)

var (
	mysqldb *MysqlDB
	once    sync.Once
)

func GetMysqlConn(MysqlDBHost string, MysqlDBPort int, MysqlDBName string, MysqlDBuser string, MysqlDBpass string) (*MysqlDB, error) {
	var err error
	once.Do(func() {
		mysqldb = &MysqlDB{MysqlDBHost, MysqlDBuser, MysqlDBName, MysqlDBpass, MysqlDBPort, nil, nil, nil}
		err = mysqldb.Conn(75, 75)
	})
	return mysqldb, err
}
