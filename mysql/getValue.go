package mysql

func (mdb *MysqlDB) getVal(query string) (string, error) {
	row := mdb.dbConn.QueryRow(query)
	var str string
	err := row.Scan(&str)
	printErrors(err)
	return str, err
}

func (mdb *MysqlDB) getValByStmt(query string, param ...interface{}) (string, error) {
	stmt, err := mdb.dbConn.Prepare(query)
	printErrors(err)
	if err != nil {
		return "", nil
	}
	defer stmt.Close()
	row := stmt.QueryRow(param...)
	var str string
	err2 := row.Scan(&str)
	printErrors(err2)
	return str, err2
}

func (mdb *MysqlDB) GetVal(qtype int, query string, param ...interface{}) (string, error) {
	lastQuery = query
	if qtype == Statement {
		return mdb.getValByStmt(query, param...)
	} else {
		return mdb.getVal(query)
	}
}
