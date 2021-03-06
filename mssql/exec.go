package mssql

func (mdb *MsSQL) stmtExec(query string, qtype int, args ...interface{}) (int64, error) {
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
		// result, err2 = rs.LastInsertId()
		result = 0
	} else if qtype == update || qtype == delete {
		result, err2 = rs.RowsAffected()
	}
	printErrors(err2)
	return result, err2
}

//Update operation ,return rows affected
func (mdb *MsSQL) Update(query string, args ...interface{}) (int64, error) {
	lastQuery = getQuery(query, args...)
	return mdb.stmtExec(query, update, args...)
}

//Insert operation ,return new insert id
func (mdb *MsSQL) Insert(query string, args ...interface{}) (int64, error) {
	lastQuery = getQuery(query, args...)
	return mdb.stmtExec(query, insert, args...)
}

//Delete operation ,return rows affected
func (mdb *MsSQL) Delete(query string, args ...interface{}) (int64, error) {
	lastQuery = getQuery(query, args...)
	return mdb.stmtExec(query, delete, args...)
}
