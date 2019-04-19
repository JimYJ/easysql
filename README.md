[![Build Status](https://travis-ci.org/JimYJ/easysql.svg?branch=master)](https://travis-ci.org/JimYJ/easysql)
[![Go Report Card](https://goreportcard.com/badge/github.com/JimYJ/easysql)](https://goreportcard.com/report/github.com/JimYJ/easysql)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/ugorji/go/master/LICENSE)
[中文说明](https://github.com/JimYJ/easysql/blob/master/README-CN.md) 

easysql encapsulated database operation, simplifying the use of the database.(include MySql,MsSQL database,more database is coming soon)

## what's new:
2.1 add mssql<br>
2.0 v2 is coming, recode tx,optimize code and structure<br>
...<br>
1.2 add cache mode<br>
1.1 add debug mode and debug func<br>
1.0 finish all base func<br>

# How to get

```
go get -u -v github.com/JimYJ/easysql/mysql
```

# Mysql driver
```go
github.com/Go-SQL-Driver/MySQL
```

# Usage

**import:**

```go
import mysql "github.com/JimYJ/easysql/mysql/v2"
```

**conn db:**
```go
mysql.Init("127.0.0.1", 3306, "dbname", "root", "123","utf8mb4", MaxIdleConns, MaxOpenConns)
mysqldb, err := mysql.Getmysqldb()//singleton pattern
or
mysqldb, err := mysql.Newmysqldb("127.0.0.1", 3306, "dbname", "root", "123", MaxIdleConns, MaxOpenConns)
```

**close conn:**
```go
mysqldb.Close()
```

**set cache timeout**
```go
mysql.SetCacheTimeout(5 * time.Second)//default 5s. if want to set cache timeout,must before trun on cache else you set timeout is no work
```
**use cache:**
```go
mysql.UseCache()
```

**close cache:**
```go
mysql.CloseCache()
```


**get value:**

```go
value,err := mysqldb.GetVal("SELECT count(*) FROM users")
or
value,err := mysqldb.GetVal("SELECT count(*) FROM users where type = ?","public")
```

**get single row data:**

```go
row,err := mysqldb.GetRow("SELECT name,email FROM users WHERE id = ?",2)
```

**get multi-rows data:**

```go
rows,err := mysqldb.GetResults("SELECT name,email FROM users where type = ?","public")
```


**If you do not want to expose the database field name,you can set field name:**
```go
mysql.SetFields([]string{"username", "useremail"})
row,err := mysqldb.GetRow("SELECT name,email FROM users WHERE id = ?",2)
or
mysql.SetFields([]string{"username", "useremail"})
rows,err := mysqldb.GetResults("SELECT name,email FROM users where type = ?","public")
```

**insert:**
```go
insertId, err := mysqldb.Insert( "insert into users set name = ?", "jim")
```

**updata:**
```go
rowsAffected, err := mysqldb.Update( "update users set name = ? where id =?", "jim", 1)
```

**delete:**
```go
rowsAffected, err := mysqldb.Delete( "delete from users where id =?", 453)
```

**transaction:**
```go
tx,_:=mysqldb.TxBegin()
insertId, err := tx.Insert("insert into users set name = ?", "jim")
rowsAffected, err := tx.Update( "update users set name = ? where id =?", "jim", 1)
rowsAffected, err := tx.Delete( "delete from users where id =?", 453)
tx.Rollback()
or
tx.Commit()
```
**debug: print last query**
```go
mysql.Debug()
```

**print all errors**
```go
mysql.DebugMode()
mysql.Init("127.0.0.1", 3306, "dbname", "root", "123","utf8mb4", 100, 100)
mysqldb, err := mysql.Getmysqldb()
...
or
mysql.DebugMode()
mysqldb, err := mysql.Newmysqldb("127.0.0.1", 3306, "dbname", "root", "123","utf8mb4", 100, 100)
```

**hide all errors**
```go
mysql.ReleaseMode()
```



