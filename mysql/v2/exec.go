package v2

func (mdb *MysqlDB) stmtExec(query string, qtype int, args ...interface{}) (int64, error) {
	stmt, err := mdb.dbConn.Prepare(query)
	printErrors(err)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
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

//Update operation ,return rows affected
func (mdb *MysqlDB) Update(query string, args ...interface{}) (int64, error) {
	lastQuery = getQuery(query, args...)
	return mdb.stmtExec(query, update, args...)
}

//Insert operation ,return new insert id
func (mdb *MysqlDB) Insert(query string, args ...interface{}) (int64, error) {
	lastQuery = getQuery(query, args...)
	return mdb.stmtExec(query, insert, args...)
}

//Delete operation ,return rows affected
func (mdb *MysqlDB) Delete(query string, args ...interface{}) (int64, error) {
	lastQuery = getQuery(query, args...)
	return mdb.stmtExec(query, delete, args...)
}
