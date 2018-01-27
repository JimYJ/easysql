package mysql

func (mdb *MysqlDB) exec(query string, qtype int, args ...interface{}) (int64, error) {
	rs, err := mdb.dbConn.Exec(query, args...)
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
	mdb.fieldlist = nil
	printErrors(err2)
	return result, err2
}

func (mdb *MysqlDB) stmtExec(query string, qtype int, args ...interface{}) (int64, error) {
	stmt, err := mdb.dbConn.Prepare(query)
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

func (mdb *MysqlDB) Update(qtype int, query string, args ...interface{}) (int64, error) {
	lastQuery = query
	if qtype == Statement {
		return mdb.stmtExec(query, update, args...)
	} else {
		return mdb.exec(query, update, args...)
	}
}

func (mdb *MysqlDB) Insert(qtype int, query string, args ...interface{}) (int64, error) {
	lastQuery = query
	if qtype == Statement {
		return mdb.stmtExec(query, insert, args...)
	} else {
		return mdb.exec(query, insert, args...)
	}
}

func (mdb *MysqlDB) Delete(qtype int, query string, args ...interface{}) (int64, error) {
	lastQuery = query
	if qtype == Statement {
		return mdb.stmtExec(query, delete, args...)
	} else {
		return mdb.exec(query, delete, args...)
	}
}
