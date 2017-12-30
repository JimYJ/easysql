package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"sync"
)

var (
	Statement     int      = 1
	Normal        int      = 0
	insert        int      = 0
	update        int      = 1
	delete        int      = 2
	customColumns []string = nil
	mysqldb       *MysqlDB
	once          sync.Once
	dBHost        string
	dBPort        int
	dBName        string
	dBuser        string
	dBpass        string
	isinit        bool
	maxIdleConns  int
	maxOpenConns  int
)

func Init(MysqlDBHost string, MysqlDBPort int, MysqlDBName string, MysqlDBuser string, MysqlDBpass string, MaxIdleConns int, MaxOpenConns int) {
	dBHost = MysqlDBHost
	dBPort = MysqlDBPort
	dBName = MysqlDBName
	dBuser = MysqlDBuser
	dBpass = MysqlDBpass
	isinit = true
	maxIdleConns = MaxIdleConns
	maxOpenConns = MaxOpenConns
}

func GetMysqlConn() (*MysqlDB, error) {
	if isinit == false {
		return nil, errors.New("DB param is not initialize!")
	}
	var err error
	once.Do(func() {
		mysqldb = &MysqlDB{dBHost, dBuser, dBName, dBpass, dBPort, nil, nil, nil}
		err = mysqldb.Conn(maxIdleConns, maxOpenConns)
	})
	return mysqldb, err
}

func NewMysqlConn(MysqlDBHost string, MysqlDBPort int, MysqlDBName string, MysqlDBuser string, MysqlDBpass string, MaxIdleConns int, MaxOpenConns int) (*MysqlDB, error) {
	var err error
	mysqldb = &MysqlDB{MysqlDBHost, MysqlDBuser, MysqlDBName, MysqlDBpass, MysqlDBPort, nil, nil, nil}
	err = mysqldb.Conn(MaxIdleConns, MaxOpenConns)
	return mysqldb, err
}

type MysqlDB struct {
	host, user, dbname, pass string
	port                     int
	dbConn                   *sql.DB
	fieldlist                []string
	tx                       *sql.Tx
}

func (self *MysqlDB) Conn(MaxIdleConns int, MaxOpenConns int) error {
	if self.host == "" || self.pass == "" || self.user == "" || self.dbname == "" {
		errs := errors.New("DB param is not initialize!")
		return errs
	}
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", self.user, self.pass, self.host, self.port, self.dbname))
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)
	err2 := db.Ping()
	if err2 != nil {
		return err2
	}
	self.dbConn = db
	return nil
}

func (self *MysqlDB) Close() {
	self.dbConn.Close()
}

func (self *MysqlDB) SetFields(fieldlist []string) {
	self.fieldlist = fieldlist
}
