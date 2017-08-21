package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"git.culiu.org/go-zk/vcutter/modules/utils/setting"
	"github.com/2liang/mcache/models/sqlmodel"
	"github.com/2liang/mcache/models/mredis"
	"strings"
	"strconv"
	"sync"
	"fmt"
	"encoding/json"
	"errors"
	"reflect"
)

type GetCacheParams struct {
	Id 		int 	`form:"id" json:"id" binding:"required"`
	Data 	string	`form:"data" json:"data" binding:"required"`
}

type DelCacheParams struct {
	Id 		int		`form:"id" json:"id" binding:"required"`
	Data 	string	`form:"data" json:"data" binding:"required"`
}

func GetCache(c *gin.Context) {
	var params GetCacheParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error("this params error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}

	// 通过ID获取相关缓存信息
	keyModel := new(sqlmodel.KeyData)
	keyModel.Id = params.Id
	keyRes, err := keyModel.GetKeyById()
	if err != nil {
		setting.SeeLog.Error("get key info error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}
	keyInfo := keyRes[0]
 	caseId := keyInfo.CaseId
	// 获取case
	caseModel := new(sqlmodel.CaseData)
	caseRes, err := caseModel.GetCaseById(caseId)
	caseInfo := caseRes[0]

	if caseInfo.Type == "REDIS" {	// redis
		res, err := RedisControl(caseInfo, keyInfo, params.Data, "GET")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": res})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "1231", "data": res})
		return
	} else if caseInfo.Type == "MC" {	// mc

	} else {
		setting.SeeLog.Error("case type " + caseInfo.Type + "is illegal")
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "case type " + caseInfo.Type + "is illegal", "data": nil})
		return
	}
}

func DeleteCache(c *gin.Context) {
	var params DelCacheParams
	if err := c.Bind(&params); err != nil {
		setting.SeeLog.Error("this params error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}

	// 通过ID获取相关缓存信息
	keyModel := new(sqlmodel.KeyData)
	keyModel.Id = params.Id
	keyRes, err := keyModel.GetKeyById()
	if err != nil {
		setting.SeeLog.Error("get key info error:" + err.Error())
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": nil})
		return
	}
	keyInfo := keyRes[0]
	caseId := keyInfo.CaseId
	// 获取case
	caseModel := new(sqlmodel.CaseData)
	caseRes, err := caseModel.GetCaseById(caseId)
	caseInfo := caseRes[0]

	if caseInfo.Type == "REDIS" {	// redis
		res, err := RedisControl(caseInfo, keyInfo, params.Data, "DEL")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": 1, "msg": err.Error(), "data": res})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "1231", "data": res})
		return
	} else if caseInfo.Type == "MC" {	// mc

	} else {
		setting.SeeLog.Error("case type " + caseInfo.Type + "is illegal")
		c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "case type " + caseInfo.Type + "is illegal", "data": nil})
		return
	}
}

func RedisControl(caseInfo sqlmodel.CaseData, keyInfo sqlmodel.KeyData, data string, method string) (interface{}, error) {
	defer func() (interface{}, error) { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err!=nil{
			setting.SeeLog.Error(err)
			return nil, errors.New("捕获到错误")
		}
		return nil, nil
	}()
	// 获取key
	var dataArr []interface{}
	err := json.Unmarshal([]byte(data), &dataArr)
	if err != nil {
		setting.SeeLog.Error("data反序列化失败")
		return nil, errors.New("data反序列化失败")
	}
	key := fmt.Sprintf(keyInfo.Prefix, dataArr[0:]...)
	redisOption := new(mredis.RedisOption)
	redisOption.Timeout 		= 3
	redisOption.WriteTimeout 	= 1
	redisOption.ReadTimeout  	= 1
	redisOption.Db				= caseInfo.Db
	if caseInfo.SlaveHost == "" {
		caseInfo.SlaveHost = caseInfo.MasterHost
	}
	if caseInfo.Port > 0 {
		// 处理配置的端口号
		caseInfo.MasterHost = caseInfo.MasterHost + ":" + strconv.Itoa(caseInfo.Port)
		shost := strings.Split(caseInfo.SlaveHost, ";")
		caseInfo.SlaveHost = ""
		temp := ""
		for i := 0; i < len(shost) - 1; i++ {
			temp = shost[i] + ":" + strconv.Itoa(caseInfo.Port)
			if i < len(shost) - 1 {
				temp += ","
			}
			caseInfo.SlaveHost += temp
		}
	}
	redisOption.MHosts = caseInfo.MasterHost
	redisOption.SHosts = caseInfo.SlaveHost
	redisClient := mredis.BaseRedis{Mutex: new(sync.Mutex)}
	redisClient.InitRedis(redisOption)
	command := ""
	if method == "GET" {
		command, err = RedisGet(keyInfo.KeyType)
	} else {
		command, err = RedisDel(keyInfo.KeyType)
	}
	setting.SeeLog.Error(command)
	if err != nil {
		return nil, err
	}
	rc := reflect.ValueOf(&redisClient)
	params := make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(key)
	refRes := rc.MethodByName(command).Call(params)
	if method == "GET" {
		return string(refRes[0].Elem().Bytes()), nil
	} else {
		setting.SeeLog.Error(refRes[0].Elem())
		if refRes[0].Elem().Int() == 1 {
			return 1, nil
		} else {
			return 0, errors.New("删除缓存数据失败~")
		}
		return refRes[0].Elem().Int(), nil
	}
}

func RedisGet(keyType string) (string, error) {
	res := ""
	switch strings.ToUpper(keyType) {
	case "STRING":
		res = "Get"
	case "HASH":
		res = "HGet"
	case "LIST":
		res = "LRANGE"
	default:
		res = ""
	}
	if res == "" {
		return "", errors.New("")
	} else {
		return res, nil
	}
}

func RedisDel(keyType string) (string, error) {
	res := ""
	switch strings.ToUpper(keyType) {
	case "STRING":
		res = "Del"
	case "HASH":
		res = "HDel"
	case "LIST":
		res = "LRANGE"
	default:
		res = ""
	}
	if res == "" {
		return "", errors.New("")
	} else {
		return res, nil
	}
}