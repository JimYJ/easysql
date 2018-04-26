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

//GetVal get single value
func (mdb *MysqlDB) GetVal(qtype int, query string, param ...interface{}) (string, error) {
	lastQuery = getQuery(query, param...)
	var rs string
	var err error
	value, found := checkCache()
	if found {
		return value.(string), nil
	}
	if qtype == Statement {
		rs, err = mdb.getValByStmt(query, param...)
		setCache(rs)
		return rs, err
	}
	rs, err = mdb.getVal(query)
	setCache(rs)
	return rs, err
}
