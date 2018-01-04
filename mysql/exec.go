package mysql

func (self *MysqlDB) exec(query string, qtype int, args ...interface{}) (int64, error) {
	rs, err := self.dbConn.Exec(query, args...)
	printErrors(err)
	if err != nil {
		return 0, err
	}
	var result int64
	var err2 error
	if qtype == insert {
		result, err2 = rs.LastInsertId()
	} else if qtype == update || qtype == delete {
		result, err2 = rs.RowsAffected()
	}
	self.fieldlist = nil
	printErrors(err2)
	return result, err2
}

func (self *MysqlDB) stmtExec(query string, qtype int, args ...interface{}) (int64, error) {
	stmt, err := self.dbConn.Prepare(query)
	printErrors(err)
	if err != nil {
		return 0, err
	}
	rs, err := stmt.Exec(args...)
	printErrors(err)
	if err != nil {
		return 0, err
	}
	var result int64
	var err2 error
	if qtype == insert {
		result, err2 = rs.LastInsertId()
	} else if qtype == update || qtype == delete {
		result, err2 = rs.RowsAffected()
	}
	printErrors(err2)
	return result, err2
}

func (self *MysqlDB) Update(qtype int, query string, args ...interface{}) (int64, error) {
	lastQuery = query
	if qtype == Statement {
		return self.stmtExec(query, update, args...)
	} else {
		return self.exec(query, update, args...)
	}
}

func (self *MysqlDB) Insert(qtype int, query string, args ...interface{}) (int64, error) {
	lastQuery = query
	if qtype == Statement {
		return self.stmtExec(query, insert, args...)
	} else {
		return self.exec(query, insert, args...)
	}
}

func (self *MysqlDB) Delete(qtype int, query string, args ...interface{}) (int64, error) {
	lastQuery = query
	if qtype == Statement {
		return self.stmtExec(query, delete, args...)
	} else {
		return self.exec(query, delete, args...)
	}
}
