//Package Dbcache database
package base

import (
	"sync"
)

var DbCache = &CacheModel{BaseModel{Mutex: new(sync.Mutex)}}

type CacheModel struct {
	BaseModel
}
