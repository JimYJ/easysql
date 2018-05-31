package mysql

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/patrickmn/go-cache"
	"log"
	"strconv"
	"strings"
	"time"
)

//Debug Print last query sql
func Debug() {
	log.Println(lastQuery)
	return
}

//DebugMode Print all errs
func DebugMode() {
	showErrors = true
}

//ReleaseMode hide all errs
func ReleaseMode() {
	showErrors = false
}

//UseCache use cache mode
func UseCache() {
	cacheMode = true
	caches = cache.New(cacheTimeout, checkCacheTimeOut)
}

//CloseCache turn off cache mode
func CloseCache() {
	cacheMode = false
	caches = nil
}

//SetCacheTimeout set cache timeout
func SetCacheTimeout(timeout time.Duration) {
	cacheTimeout = timeout
}

func printErrors(err error) {
	if err != nil && showErrors == true {
		log.Println(err)
	}
}

func getQuery(query string, param ...interface{}) string {
	if param == nil {
		return query
	}
	queryFormat := strings.Replace(query, "?", "%v", -1)
	return fmt.Sprintf(queryFormat, param...)
}

func hashsha1(data string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(data))
	return hex.EncodeToString(sha1.Sum([]byte("")))
}

func checkCache() (interface{}, bool) {
	if cacheMode {
		key := hashsha1(lastQuery)
		return caches.Get(key)
	}
	return nil, false
}
func setCache(value interface{}) {
	if cacheMode {
		key := hashsha1(lastQuery)
		caches.Set(key, value, cacheTimeout)
	}
}

func anyToString(m interface{}) string {
	switch m.(type) {
	case string:
		return m.(string)
	case int:
		return strconv.Itoa(m.(int))
	case int64:
		return strconv.FormatInt(m.(int64), 10)
	default:
		return ""
		//处理业务
	}
}
