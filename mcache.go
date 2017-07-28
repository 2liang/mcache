//Goil powerd chuchujie group
package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/2liang/mcache/models"
	"github.com/2liang/mcache/modules/utils/setting"
	"github.com/2liang/mcache/routers"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	setting.Daemon()

	setting.HookReload = make([]func(), 0)
	// init models
	setting.HookReload = append(setting.HookReload, models.Init)
	setting.LoadConfig()

	r := gin.New()
	buffer_size := 20
	if setting.IsProMode {
		gin.SetMode(gin.ReleaseMode)
		buffer_size = 51200
	}

	accessLogger := &setting.Logwriter{Mutex: new(sync.Mutex), Prefix: "access", BufferSize: buffer_size}
	defer func() {
		accessLogger.Flush()
	}()

	r.Use(gin.LoggerWithWriter(accessLogger), gin.RecoveryWithWriter(setting.Logwriterfile))

	serve := &http.Server{
		Addr:         setting.AppHost,
		Handler:      r,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 10,
	}
	routers.Init(r)

	setting.SetSignal(serve)
}
