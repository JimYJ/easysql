[![Build Status](https://travis-ci.org/JimYJ/easysql.svg?branch=master)](https://travis-ci.org/JimYJ/easysql)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/ugorji/go/master/LICENSE)
[中文说明](https://github.com/JimYJ/easysql/blob/master/README-CN.md) 

easysql encapsulated database operation, simplifying the use of the database.(include Mysql database,more database is coming soon)

## what's up:
1.2 add cache mode
1.1 add debug mode and debug func
1.0 finish all base func

# How to get

```
go get github.com/JimYJ/easysql/mysql
```

# Mysql driver
```go
github.com/Go-SQL-Driver/MySQL
```

# Usage

**import:**

```go
import "github.com/JimYJ/easysql/mysql"
```

**conn db:**
```go
mysql.Init("127.0.0.1", 3306, "dbname", "root", "123", MaxIdleConns, MaxOpenConns)
mysqldb, err := mysql.Getmysqldb()//singleton pattern
or
mysqldb, err := mysql.Newmysqldb("127.0.0.1", 3306, "dbname", "root", "123", MaxIdleConns, MaxOpenConns)
```

**close conn:**
```go
mysqldb.Close()
```
**use cache:**
```go
mysql.UseCache()
```

**close cache:**
```go
mysql.CloseCache()
```

**set cache timeout**
```go
mysql.SetCacheTimeout(5 * time.Second)
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
mysqldb.TxBegin()
insertId, err := mysqldb.TxInsert(mysql.Normal, "insert into users set name = ?", "jim")
rowsAffected, err := mysqldb.TxUpdate(mysql.Normal, "update users set name = ? where id =?", "jim", 1)
rowsAffected, err := mysqldb.TxDelete(mysql.Normal, "delete from users where id =?", 453)
mysqldb.TxRollback()
or
mysqldb.TxCommit()
```

**transaction with statement:**
```go
mysqldb.TxBegin()
insertId, err := mysqldb.TxInsert(mysql.Statement, "insert into users set name = ?", "jim")
rowsAffected, err := mysqldb.TxUpdate(mysql.Statement, "update users set name = ? where id =?", "jim", 1)
rowsAffected, err := mysqldb.TxDelete(mysql.Statement, "delete from users where id =?", 453)
mysqldb.TxRollback()
or
mysqldb.TxCommit()
```

**debug: print last query**
```go
mysql.Debug()
```

**print all errors**
```go
mysql.ShowErrors()
mysql.Init("127.0.0.1", 3306, "dbname", "root", "123", 100, 100)
mysqldb, err := mysql.Getmysqldb()
...
or
mysql.ShowErrors()
mysqldb, err := mysql.Newmysqldb("127.0.0.1", 3306, "dbname", "root", "123", 100, 100)
```

**hide all errors**
```go
mysql.HideErrors()
```



