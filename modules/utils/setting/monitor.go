/*
作者：陆恒
邮箱：luheng@chuchujie.com
时间：2017/3/20
功能：上报
*/

package setting

import (
	"net"
	"strconv"
	"time"
)

/*
  #Created by Luheng on 2017/3/20.
  #Arguments: 编号, 次数
  #Return:
  #Description: 上报访问次数
*/
func SendCount(pid int, cnt int) {

	if !IsProMode {
		return
	}
	conn, err := net.Dial("udp", "10.30.10.227:21555")
	defer conn.Close()
	if err != nil {
		SeeLog.Error("监控UDP失败", err)
	}
	msg := "go-st" + "-==-" + time.Now().Format("2006-01-02 15:04:05") + "-==-" + strconv.Itoa(pid) + "-==-" + strconv.Itoa(cnt)
	conn.Write([]byte(msg))
	SeeLog.Info("上报日志", msg)
}

/*
  #Created by Luheng on 2017/3/20.
  #Arguments: 编号, 时间(ms)
  #Return:
  #Description: 上报延迟接口
*/
func SendTime(pid int, timed int) {

	if !IsProMode {
		return
	}
	conn, err := net.Dial("udp", "10.30.10.227:21556")
	defer conn.Close()
	if err != nil {
		SeeLog.Error("监控UDP失败", err)
	}
	msg := "go-st" + "-==-" + time.Now().Format("2006-01-02 15:04:05") + "-==-" + strconv.Itoa(pid) + "-==-" + strconv.Itoa(timed)
	conn.Write([]byte(msg))
	//SeeLog.Info("上报日志", msg)
}
