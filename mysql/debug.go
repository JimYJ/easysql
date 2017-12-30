package mysql

func (self *MysqlDB) GetLastQuery() string {
	return lastQuery
}
