package mysql

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/patrickmn/go-cache"
	"log"
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
}

//OffCache turn off cache mode
func CloseCache() {
	cacheMode = false
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
		c := cache.New(cacheTimeout, checkCacheTimeOut)
		value, found := c.Get(key)
		if found {
			return value, true
		}
		return nil, false
	}
	return nil, false
}
func setCache(value interface{}) {
	if cacheMode {
		key := hashsha1(lastQuery)
		c := cache.New(cacheTimeout, checkCacheTimeOut)
		c.Set(key, value, cacheTimeout)
	}
}
