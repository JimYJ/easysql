package mysql

func (self *MysqlDB) exec(query string, qtype int, args ...interface{}) (int64, error) {
	rs, err := self.dbConn.Exec(query, args...)
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
	return result, err2
}

func (self *MysqlDB) stmtExec(query string, qtype int, args ...interface{}) (int64, error) {
	stmt, err := self.dbConn.Prepare(query)
	if err != nil {
		return 0, err
	}
	rs, err := stmt.Exec(args...)
	if err != nil {
		if qtype == insert {
			return 0, err
		} else if qtype == update {
			return 0, err
		} else if qtype == delete {
			return 0, err
		}
	}
	var result int64
	var err2 error
	if qtype == insert {
		result, err2 = rs.LastInsertId()
	} else if qtype == update || qtype == delete {
		result, err2 = rs.RowsAffected()
	}
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
