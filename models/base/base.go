//Package base
package base

import (
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"

	"github.com/2liang/mcache/modules/utils/setting"
)

//BaseModel
type BaseModel struct {
	*sync.Mutex
	m  []*xorm.Engine
	s  []*xorm.Engine
	c  int32
	cc int32
}

//BaseModel.GetMaster
func (base *BaseModel) GetMaster() *xorm.Engine {

	base.setPool(false, nil)
	if err := base.m[0].Ping(); err != nil {
		base.m[0].Logger().Error("the DB :" + base.m[0].DataSourceName() + " is offline")
	}
	//fmt.Printf("%+v@%+v", base.m[0], len(base.m))
	return base.m[0]
}

//BaseModel.GetSlave
func (base *BaseModel) GetSlave() *xorm.Engine {

	base.setPool(false, nil)

	atomic.AddInt32(&base.cc, 1)
	atomic.CompareAndSwapInt32(&base.cc, 1e9, 0)
	d := base.s[(base.cc % int32(len(base.s)))]
	if err := d.Ping(); err != nil {
		d.Logger().Error("the DB :" + d.DataSourceName() + " is offline")
		atomic.AddInt32(&base.c, 1)
		if atomic.CompareAndSwapInt32(&base.c, int32(len(base.s)), 0) {
			goto HEAR
		}
		return base.GetSlave()
	} else {
		return d
	}

HEAR:
	return d
}

//BaseModel.SetPool
func (base *BaseModel) setPool(reset bool, resetMap map[string][]*xorm.Engine) {
	defer base.Mutex.Unlock()
	base.Mutex.Lock()

	if reset {
		tmp_m := base.m
		tmp_s := base.s
		base.m = resetMap["m"]
		base.s = resetMap["s"]
		base.Close(tmp_m)
		base.Close(tmp_s)
	}

}

//BaseModel.Close
func (base *BaseModel) Close(x []*xorm.Engine) {
	if len(x) > 0 {
		time.AfterFunc(time.Second, func() {
			for _, _x := range x {
				//fmt.Printf("%+v****************\n", x)
				_x.Close()
			}
		})
	}
}

//BaseModel.InitXorm init the xorm group
func (base *BaseModel) InitXorm(group string) {
	mm := make(map[string]interface{}, 5)

	mm["source"] = setting.Cfg.MustValue(group+".master", "#1", "")

	if len(mm["source"].(string)) < 1 {
		setting.Logger.Panic("The dbs [" + group + "] master init error")
	}

	mm["driver"] = "mysql"
	mm["maxidle"] = setting.Cfg.MustInt(group+".master", "max_idle_conn", 30)
	mm["maxopen"] = setting.Cfg.MustInt(group+".master", "max_open_conn", 50)
	mm["debug"] = setting.Cfg.MustBool(group+".master", "debug_log", true)
	m := make([]*xorm.Engine, 0)
	m = append(m, base.InitDb(mm))

	dsns := setting.Cfg.MustValue(group+".slave", "#1", "")
	if len(dsns) < 1 {
		setting.Logger.Panic("The dbs [" + group + "] slave init error")
	}

	dsn := strings.Split(dsns, ",")
	if len(dsn) < 1 {
		setting.Logger.Panic("The dbs [" + group + "] slave init error")
	}

	mm["maxidle"] = setting.Cfg.MustInt(group+".slave", "max_idle_conn", 30)
	mm["maxopen"] = setting.Cfg.MustInt(group+".slave", "max_open_conn", 50)
	mm["debug"] = setting.Cfg.MustBool(group+".slave", "debug_log", true)

	s := make([]*xorm.Engine, 0)
	for _, d := range dsn {
		mm["source"] = d
		//cacher := xorm.NewLRUCacher2(xorm.NewMemoryStore(), 1800*time.Second, 1000)
		engine := base.InitDb(mm)
		//engine.SetDefaultCacher(cacher)
		s = append(s, engine)
	}

	reset := make(map[string][]*xorm.Engine, 2)
	reset["m"] = m
	reset["s"] = s
	base.setPool(true, reset)
}

//BaseModel.initDb
func (base *BaseModel) InitDb(m map[string]interface{}) *xorm.Engine {

	var err error
	var orm *xorm.Engine
	orm, err = xorm.NewEngine(m["driver"].(string), m["source"].(string))
	if err != nil {
		panic(err)
	}

	os.MkdirAll("./logs", os.ModePerm)
	f, err := os.OpenFile("logs/sql.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0664)
	if err != nil {
		setting.Logger.Panic("create dblog file failed:", err)
	}

	orm.SetMaxIdleConns(m["maxidle"].(int))
	orm.SetMaxOpenConns(m["maxopen"].(int))

	if !setting.IsProMode {
		orm.ShowSQL(true)
		orm.Logger().SetLevel(core.LOG_DEBUG)
	} else {
		orm.ShowSQL(false)
		orm.Logger().SetLevel(core.LOG_WARNING)
	}

	orm.SetLogger(xorm.NewSimpleLogger(f))

	return orm
}
