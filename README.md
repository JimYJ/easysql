[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/ugorji/go/master/LICENSE)

EasyDB encapsulated database operation, simplifying the use of the database.(include Mysql database,more database is coming soon)

# How to get

```
go get github.com/JimYJ/EasyDB/mysql
```

# Usage

**import:**

```go
import "github.com/JimYJ/EasyDB/mysql"
```

**conn db:**
```go
mysql.Init("127.0.0.1", 3306, "dbname", "root", "123", 100, 100)
mysqldb, err := mysql.GetMysqlConn()//singleton pattern
or
mysqldb, err := mysql.NewMysqlConn("127.0.0.1", 3306, "dbname", "root", "123", 100, 100)
```

**get value:**

```go
value,err := mysqldb.GetVal(mysql.Normal,"SELECT count(*) FROM users")
```
**get value with statement:**

```go
value,err := mysqldb.GetVal(mysql.Statement,"SELECT count(*) FROM users")
or
value,err := mysqldb.GetVal(mysql.Statement,"SELECT count(*) FROM users where type = ?","public")
```

**get single row data:**
```go
row,err := mysqldb.GetRow(mysql.Normal,"SELECT name,email FROM users WHERE id = 2")
```

**get single row data with statement:**
```go
row,err := mysqldb.GetRow(mysql.Statement,"SELECT name,email FROM users WHERE id = ?",2)
```

**get multi-rows data:**
```go
rows,err := mysqldb.GetResults(mysql.Normal,"SELECT name,email FROM users where type = 'public'")
```

**get multi-rows data with statement:**
```go
rows,err := mysqldb.GetResults(mysql.Statement,"SELECT name,email FROM users where type = ?","public")
```


**If you do not want to expose the database field name,you can set field name:**
```go
mysql.SetFields([]string{"username", "useremail"})
row,err := mysqldb.GetRow(mysql.Statement,"SELECT name,email FROM users WHERE id = ?",2)
or
mysql.SetFields([]string{"username", "useremail"})
rows,err := mysqldb.GetResults(mysql.Statement,"SELECT name,email FROM users where type = ?","public")
```

**insert:**
```go
insertId, err := mysqldb.Insert(mysql.Normal, "insert into users set name = ?", "jim")
```


**insert with statement:**
```go
insertId, err := mysqldb.Insert(mysql.Statement, "insert into users set name = ?", "jim")
```

**updata:**
```go
rowsAffected, err := mysqldb.Update(mysql.Normal, "update users set name = ? where id =?", "jim", 1)
```

**updata with statement:**
```go
rowsAffected, err := mysqldb.Update(mysql.Statement, "update users set name = ? where id =?", "jim", 1)
```

**delete:**
```go
rowsAffected, err := mysqldb.Delete(mysql.Normal, "delete from users where id =?", 453)
```

**delete with statement:**
```go
rowsAffected, err := mysqldb.Delete(mysql.Statement, "delete from users where id =?", 453)
```

**transaction:**
```go
mysqlconn.TxBegin()
insertId, err := mysqlconn.TxInsert(mysql.Normal, "insert into users set name = ?", "jim")
rowsAffected, err := mysqlconn.TxUpdate(mysql.Normal, "update users set name = ? where id =?", "jim", 1)
rowsAffected, err := mysqlconn.TxDelete(mysql.Normal, "delete from users where id =?", 453)

mysqlconn.TxRollback()
or
mysqlconn.TxCommit()
```

**transaction with statement:**
```go
mysqlconn.TxBegin()
insertId, err := mysqlconn.TxInsert(mysql.Statement, "insert into users set name = ?", "jim")
rowsAffected, err := mysqlconn.TxUpdate(mysql.Statement, "update users set name = ? where id =?", "jim", 1)
rowsAffected, err := mysqlconn.TxDelete(mysql.Statement, "delete from users where id =?", 453)

mysqlconn.TxRollback()
or
mysqlconn.TxCommit()
```



