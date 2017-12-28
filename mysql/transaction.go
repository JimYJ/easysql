package mysql

import (
	"errors"
)

func (self *MysqlDB) TxBegin() error {
	var err error
	self.tx, err = self.dbConn.Begin()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (self *MysqlDB) TxCommit() error {
	if self.tx == nil {
		err := errors.New("Transaction didn't initializtion!")
		return err
	}
	err := self.tx.Commit()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (self *MysqlDB) TxRollback() error {
	if self.tx == nil {
		err := errors.New("Transaction didn't initializtion!")
		return err
	}
	err := self.tx.Rollback()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (self *MysqlDB) txExec(query string, qtype int, args ...interface{}) (int64, error) {
	if self.tx == nil {
		err := errors.New("Transaction didn't initializtion!")
		return 0, err
	}
	rs, err := self.tx.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	var result int64
	var err2 error
	if qtype == Insert {
		result, err2 = rs.LastInsertId()
	} else if qtype == Update {
		result, err2 = rs.RowsAffected()
	}
	self.fieldlist = nil
	return result, err2
}

func (self *MysqlDB) stmtTxExec(query string, qtype int, args ...interface{}) (int64, error) {
	if self.tx == nil {
		err := errors.New("Transaction didn't initializtion!")
		return 0, err
	}
	stmt, err := self.tx.Prepare(query)
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

func (self *MysqlDB) TxUpdate(qtype int, query string, args ...interface{}) (int64, error) {
	if qtype == Statement {
		return self.stmtTxExec(query, Update, args...)
	} else {
		return self.txExec(query, Update, args...)
	}
}

func (self *MysqlDB) TxInsert(qtype int, query string, args ...interface{}) (int64, error) {
	if qtype == Statement {
		return self.stmtTxExec(query, Insert, args...)
	} else {
		return self.txExec(query, Insert, args...)
	}
}

func (self *MysqlDB) TxDelete(qtype int, query string, args ...interface{}) (int64, error) {
	if qtype == Statement {
		return self.stmtTxExec(query, Delete, args...)
	} else {
		return self.txExec(query, Delete, args...)
	}
}
