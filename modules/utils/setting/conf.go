//Package setting provide goil's settings
package setting

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/howeyc/fsnotify"
)

const (
	APP_VER = "0.10"
)

var (
	AppName   string
	AppHost   string
	AppVer    string
	IsProMode bool
	TimeZone  string
)

var (
	Cfg *goconfig.ConfigFile
)

var (
	AppConfPath          = "conf/app.ini"
	HookReload           []func()
	FigureSite           string
	VideoSite            string
	AwsBucket            string
	AwsRegion            string
	AwsAccessKeyId       string
	AwsSecretAccessKey   string
	QiniuAccessKeyId     string
	QiniuSecretAccessKey string
	CacheNodes           []string
	QiniuBucket          string
	QiniuVmp4Prefix      string
	QiniuVm3u8Prefix     string
	QiniuPicPrefix       string
)

// LoadConfig loads configuration file.
func LoadConfig() *goconfig.ConfigFile {
	var err error

	if fh, _ := os.OpenFile(AppConfPath, os.O_RDONLY|os.O_CREATE, 0600); fh != nil {
		fh.Close()
	}

	// Load configuration, set app version and Log level.
	Cfg, err = goconfig.LoadConfigFile(AppConfPath)

	if err != nil {
		Logger.Error("Fail to load configuration file: " + err.Error())
		os.Exit(2)
	}

	//Cfg.BlockMode = false

	// set time zone of wetalk system
	TimeZone = Cfg.MustValue("app", "time_zone", "UTC")
	if _, err := time.LoadLocation(TimeZone); err == nil {
		os.Setenv("TZ", TimeZone)
	} else {
		Logger.Error("Wrong time_zone: " + TimeZone + " " + err.Error())
		os.Exit(2)
	}

	//aws config
	FigureSite = Cfg.MustValue("site", "figure_site", "")
	VideoSite = Cfg.MustValue("site", "video_site", "")

	AwsBucket = Cfg.MustValue("aws", "aws_bucket")
	AwsRegion = Cfg.MustValue("aws", "aws_region")
	AwsAccessKeyId = Cfg.MustValue("aws", "aws_access_key_id")
	AwsSecretAccessKey = Cfg.MustValue("aws", "aws_secret_access_key")

	QiniuAccessKeyId = Cfg.MustValue("qiniu", "qiniu_access_key_id")
	QiniuSecretAccessKey = Cfg.MustValue("qiniu", "qiniu_secret_access_key")
	QiniuBucket = Cfg.MustValue("qiniu", "bucket")
	QiniuVmp4Prefix = Cfg.MustValue("qiniu", "vmp4_prefix")
	QiniuVm3u8Prefix = Cfg.MustValue("qiniu", "vm3u8_prefix")
	QiniuPicPrefix = Cfg.MustValue("qiniu", "pic_prefix")

	os.MkdirAll("./tmpcache", os.ModePerm)

	// Trim 4th part.
	AppVer = strings.Join(strings.Split(APP_VER, ".")[:2], ".")
	AppHost = Cfg.MustValue("app", "app_host", "9501")

	IsProMode = Cfg.MustValue("app", "run_mode") == "pro"

	if IsDebug {
		IsProMode = !IsDebug
	}

	reloadConfig()
	configWatcher()

	return Cfg
}

func reloadConfig() {
	AppName = Cfg.MustValue("app", "app_name", "WeTalk Community")
	CacheNodes = strings.Split(Cfg.MustValue("cachenodes", "nodes", ""), ",")

	for _, f := range HookReload {
		f()
	}
}

var eventTime = make(map[string]int64)

func configWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic("Failed start app watcher: " + err.Error())
	}

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				switch filepath.Ext(event.Name) {
				case ".ini":
					if checkEventTime(event.Name) {
						continue
					}

					Logger.Info(event)
					if err := Cfg.Reload(); err != nil {
						Logger.Error("Conf Reload: ", err)
					}

					reloadConfig()
					Logger.Info("Config Reloaded")
				}
			}
		}
	}()

	if err := watcher.WatchFlags("conf", fsnotify.FSN_MODIFY); err != nil {
		Logger.Error(err)
	}
}

// checkEventTime returns true if FileModTime does not change.
func checkEventTime(name string) bool {
	mt := getFileModTime(name)
	if eventTime[name] == mt {
		return true
	}

	eventTime[name] = mt
	return false
}

// getFileModTime retuens unix timestamp of `os.File.ModTime` by given path.
func getFileModTime(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {
		Logger.Error("Fail to open file[ %s ]\n", err)
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		Logger.Error("Fail to get file information[ %s ]\n", err)
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}
