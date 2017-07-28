//Package models provide the basic model of the db
package models

import (
	"github.com/2liang/mcache/models/base"
)

func Init() {
	base.DbCache.InitXorm("dbcache")
	base.DbCache.GetMaster()
	base.DbCache.GetSlave()
}
