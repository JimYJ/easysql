package mysql

func (self *MysqlDB) exec(query string, qtype int, args ...interface{}) (int64, error) {
	rs, err := self.dbConn.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	var result int64
	var err2 error
	if qtype == Insert {
		result, err2 = rs.LastInsertId()
	} else if qtype == Update || qtype == Delete {
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
		if qtype == Insert {
			return 0, err
		} else if qtype == Update {
			return 0, err
		} else if qtype == Delete {
			return 0, err
		}
	}
	var result int64
	var err2 error
	if qtype == Insert {
		result, err2 = rs.LastInsertId()
	} else if qtype == Update || qtype == Delete {
		result, err2 = rs.RowsAffected()
	}
	return result, err2
}

func (self *MysqlDB) Update(qtype int, query string, args ...interface{}) (int64, error) {
	if qtype == Statement {
		return self.stmtExec(query, Update, args...)
	} else {
		return self.exec(query, Update, args...)
	}
}

func (self *MysqlDB) Insert(qtype int, query string, args ...interface{}) (int64, error) {
	if qtype == Statement {
		return self.stmtExec(query, Insert, args...)
	} else {
		return self.exec(query, Insert, args...)
	}
}

func (self *MysqlDB) Delete(qtype int, query string, args ...interface{}) (int64, error) {
	if qtype == Statement {
		return self.stmtExec(query, Delete, args...)
	} else {
		return self.exec(query, Delete, args...)
	}
}
