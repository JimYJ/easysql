package mssql

func (mdb *MsSQL) getValByStmt(query string, param ...interface{}) (interface{}, error) {
	stmt, err := mdb.dbConn.Prepare(query)
	printErrors(err)
	if err != nil {
		return "", nil
	}
	defer stmt.Close()
	row := stmt.QueryRow(param...)
	var value interface{}
	err2 := row.Scan(&value)
	printErrors(err2)
	return value, err2
}

//GetVal get single value
func (mdb *MsSQL) GetVal(query string, param ...interface{}) (interface{}, error) {
	lastQuery = getQuery(query, param...)
	var value interface{}
	var err error
	value, found := checkCache()
	if found {
		return value, nil
	}
	value, err = mdb.getValByStmt(query, param...)
	b, ok := value.([]byte)
	if ok {
		value = string(b)
	}
	setCache(value)
	return value, err
}
