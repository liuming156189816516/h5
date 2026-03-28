package utils

import (
	"strings"
	"time"
)
import "fmt"

const (
	TimeFormat_Full           = "2006-01-02 15:04:05.999"
	TimeFormat_Sec            = "2006-01-02 15:04:05"
	TimeFormat_Sec_New        = "2006/01/02 15:04:05"
	TimeFormat_Sec_NoZero     = "2006-1-2 15:4:5"
	TimeFormat_Minute         = "2006-01-02 15:04"
	TimeFormat_Hour           = "2006-01-02 15"
	TimeFormat_Date           = "2006-01-02"
	TimeFormat_IDate          = "20060102"
	TimeFormat_Month          = "200601"
	TimeFormat_LianTime       = "20060102150405"
	TimeFormat_HourTime       = "2006010215"
	TimeFormat_DateTime       = "01-02"
	TimeFormat_HourMinute     = "15:04"
	TimeFormat_HourMinuteTime = "15:04:05"
	TimeZeroPointsDff         = 72000 //与北京时间相差
)

/**
CST时间转到美东时间 utc-4
*/
func TimeToZoneUTC4(str string) (string, error) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return "", err
	}
	local, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(TimeFormat_Sec, str, local)
	return theTime.In(loc).Format(TimeFormat_Sec), nil
}

// 返回 "15:50"
func GetHourMin() string {
	return time.Now().Format(TimeFormat_HourMinute)
}

func GetYMD() string {
	return time.Now().Format(TimeFormat_IDate)
}

//月日
func GetMd(t int64) string {
	return time.Unix(t, 0).Format(TimeFormat_DateTime)
}

// 年月日
func GetYmd(t int64) string {
	return time.Unix(t, 0).Format(TimeFormat_IDate)
}
//获取年月日
func GetYmd2(t int64) string {
	return time.Unix(t, 0).Format(TimeFormat_Date)
}



func GetDay(t int64) string {
	str := time.Unix(t, 0).Format(TimeFormat_Date)
	s := strings.Split(str, "-")
	if len(s) < 3 {
		return ""
	}
	return s[2]
}

//获取现在的时间
func GetNowTimeString() string {
	return time.Now().Format(TimeFormat_Sec)
}

//时间搓转时间字符串
func TimeIntToString(t int64) string {
	return time.Unix(t, 0).Format(TimeFormat_Sec)
}

//时间搓转时间字符串
func TimeIntToLianString(t int64) string {
	return time.Unix(t, 0).Format(TimeFormat_LianTime)
}

//时间搓转时间字符串
func ToStringTime(t int64) string {
	if t == 0 {
		return ""
	}
	return time.Unix(t, 0).Format(TimeFormat_Sec)
}

//时间搓转时间字符串
func TimeStrToTime(str string) int64 {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(TimeFormat_Sec, str, loc)
	return theTime.Unix()
}

func TimeStrYearToTime(t time.Time, str string) int64 {
	str = fmt.Sprintf("%d-%s", t.Year(), str)
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(TimeFormat_Sec_NoZero, str, loc)
	return theTime.Unix()
}

//获取 月
func GetIMonth(t int64) int64 {
	return StrToInt64(time.Unix(t, 0).Format(TimeFormat_Month))
}

func ParseMonthToTime(v int64) time.Time {
	t, _ := time.Parse(TimeFormat_Month, IntToStr(v))
	return t
}

//获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

//获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

//获取idate
func GetIDate(t int64) int64 {
	return StrToInt64(time.Unix(t, 0).Format(TimeFormat_IDate))
}

/**
获取int64时间戳 天
*/
func GetTimeDay(t int64) int {
	return time.Unix(t, 0).Day()
}

//获取idate
func GetStrIDate(t int64) string {
	return time.Unix(t, 0).Format(TimeFormat_IDate)
}

func GetStrDate(t int64) string {
	return time.Unix(t, 0).Format(TimeFormat_Date)
}

//str idata + n
func StrIDateInc(idate string, n int64) string {

	l := len(idate) //不检查 idate
	if l < 6 {
		return "20060102"
	}
	newDate := idate[0:l-4] + "-" + idate[l-4:l-2] + "-" + idate[l-2:l] + " 00:00:00"
	newIdate := TimeStrToTime(newDate) + n*(24*3600)
	return GetStrIDate(newIdate)
}

//int idata + n
func IDateInc(idate int64, n int64) int64 {
	return StrToInt64(StrIDateInc(IntToStr(idate), n))
}

//获取idate
func GetNowStrIDate(nows ...time.Time) string {
	var now time.Time
	if len(nows) > 0 {
		now = nows[0]
	} else {
		now = time.Now()
	}
	return now.Format(TimeFormat_IDate)
}

//获取现在到 今天晚点0点的时间差
func GetTomorryZeroTimeDiff(nows ...time.Time) int64 {
	var now time.Time
	if len(nows) > 0 {
		now = nows[0]
	} else {
		now = time.Now()
	}
	timeStr := now.Format(TimeFormat_Date)
	t, _ := time.ParseInLocation(TimeFormat_Sec, timeStr+" 23:59:59", time.Local)
	return (t.Unix() + 1 - now.Unix())
}

//获取现在到 今天凌晨4点的时间差
func GetTomorryFourTimeDiff(nows ...time.Time) int64 {
	var now time.Time
	if len(nows) > 0 {
		now = nows[0]
	} else {
		now = time.Now()
	}
	timeStr := now.Format(TimeFormat_Date)
	t, _ := time.ParseInLocation(TimeFormat_Sec, timeStr+" 03:59:59", time.Local)
	return (t.Unix() + 1 - now.Unix())
}

//明天0点
func GetTomorryZeroTime(nows ...time.Time) int64 {
	var now time.Time
	if len(nows) > 0 {
		now = nows[0]
	} else {
		now = time.Now()
	}
	timeStr := now.Format(TimeFormat_Date)
	t, _ := time.ParseInLocation(TimeFormat_Sec, timeStr+" 23:59:59", time.Local)
	return t.Unix() + 1
}

//今天0点
func GetDayZeroTime(nows ...time.Time) int64 {
	var now time.Time
	if len(nows) > 0 {
		now = nows[0]
	} else {
		now = time.Now()
	}
	timeStr := now.Format(TimeFormat_Date)
	t, _ := time.ParseInLocation(TimeFormat_Sec, timeStr+" 00:00:00", time.Local)
	return t.Unix()
}

//根据2006-01-01格式的日期，获取相应的0点
func GetDayZeroTimeByTimeStr(timeStr string) int64 {
	t, _ := time.ParseInLocation(TimeFormat_Sec, timeStr+" 00:00:00", time.Local)
	return t.Unix()
}

//今天0点
func GetDayZero(t int64) int64 {
	timeStr := time.Unix(t, 0).Format(TimeFormat_Date)
	tt, _ := time.ParseInLocation(TimeFormat_Sec, timeStr+" 00:00:00", time.Local)
	return tt.Unix()
}

// 获取45天前0点
func GetDataStartTime() int64 {
	day_start := GetDayZeroTime()
	return day_start - (45 * 24 * 3600)
}

// 获取45天前0点 str
func GetDataStartStrTime() string {
	return TimeIntToString(GetDataStartTime())
}

// 本周的零点
func GetWeekStartStr() string {
	return TimeIntToString(GetWeekStart())
}

// 下周的零点
func GetNextWeekStartStr() int64 {
	return GetWeekStart() + 86400*7
}

// 下一月
func GetNextMonthZero() int64 {
	i := time.Now().Unix() + 30*86400
	return GetMonthStartByTime(i)
}

// 本月的零点
func GetMonthStart() int64 {
	month := GetIDate(GetDayZeroTime()) / 100      //本月
	return GetTimeByIdate(IntToStr(month*100 + 1)) //本月第一秒
}

// 某时间段月零点
func GetMonthStartByTime(t int64) int64 {
	month := GetIDate(t) / 100                     //本月
	return GetTimeByIdate(IntToStr(month*100 + 1)) //本月第一秒
}

// 下一月 0点
func GetNextMonthZeroByItime(itime int64) int64 {
	r := GetIDate(itime) % 100 //日
	inc := int64(31)
	if r > 5 {
		r = 25
	}
	month := GetIDate(itime+inc*24*3600) / 100
	return GetTimeByIdate(IntToStr(month*100 + 1)) //本月第一秒
}

// 某小时起始
func GetHourStartByTime(t int64) int64 {
	timeStr := time.Unix(t, 0).Format(TimeFormat_Hour)
	tt, _ := time.ParseInLocation(TimeFormat_Sec, timeStr+":00:00", time.Local)
	return tt.Unix()
}

// 本周的零点int
func GetWeekStart() int64 {
	todayStart := GetDayZeroTime()
	weekDay := time.Now().Weekday()
	i := 0
	switch weekDay {
	case time.Monday:
		{
			i = 0
		}
	case time.Tuesday:
		{
			i = 1
		}
	case time.Wednesday:
		{
			i = 2
		}
	case time.Thursday:
		{
			i = 3
		}
	case time.Friday:
		{
			i = 4
		}
	case time.Saturday:
		{
			i = 5
		}
	case time.Sunday:
		{
			i = 6
		}
	}
	return todayStart - int64(i*3600*24)
}

// // 所穿时间 周 1 itime
func GetWeekStartItime(itime int64) int64 {
	zore := GetDayZero(itime)
	weekDay := time.Unix(zore, 0).Weekday()
	i := 0
	switch weekDay {
	case time.Monday:
		{
			i = 0
		}
	case time.Tuesday:
		{
			i = 1
		}
	case time.Wednesday:
		{
			i = 2
		}
	case time.Thursday:
		{
			i = 3
		}
	case time.Friday:
		{
			i = 4
		}
	case time.Saturday:
		{
			i = 5
		}
	case time.Sunday:
		{
			i = 6
		}
	}
	tmp := zore - int64(i*3600*24)
	return tmp
}

// // 所穿时间 周 1 idate
func GetWeekStartByIdate(idate int64) int64 {
	r := idate % 100
	y := (idate % 10000) / 100
	n := idate / 10000
	zore := TimeStrToTime(fmt.Sprintf("%d-%02d-%02d 00:00:00", n, y, r))
	weekDay := time.Unix(zore, 0).Weekday()
	i := 0
	switch weekDay {
	case time.Monday:
		{
			i = 0
		}
	case time.Tuesday:
		{
			i = 1
		}
	case time.Wednesday:
		{
			i = 2
		}
	case time.Thursday:
		{
			i = 3
		}
	case time.Friday:
		{
			i = 4
		}
	case time.Saturday:
		{
			i = 5
		}
	case time.Sunday:
		{
			i = 6
		}
	}
	tmp := zore - int64(i*3600*24)
	return GetIDate(tmp)
}

//idate 星期几
func GetWeekByIdate(idate int64) int32 {
	r := idate % 100
	y := (idate % 10000) / 100
	n := idate / 10000
	zore := TimeStrToTime(fmt.Sprintf("%d-%02d-%02d 00:00:00", n, y, r))
	weekDay := time.Unix(zore, 0).Weekday()
	return int32(weekDay)
}

// 本周的零点int
func GetWeekByStrIdate(idate string) string {
	l := len(idate) //不检查 idate
	newDate := TimeStrToTime(idate[0:l-4] + "-" + idate[l-4:l-2] + "-" + idate[l-2:l] + " 00:00:00")
	weekDay := time.Unix(newDate, 0).Weekday()

	switch weekDay {
	case time.Sunday:
		return "日"
	case time.Monday:
		return "一"
	case time.Tuesday:
		return "二"
	case time.Wednesday:
		return "三"
	case time.Thursday:
		return "四"
	case time.Friday:
		return "五"
	case time.Saturday:
		return "六"
	}
	return ""
}

//idate 转 0点时间戳
func GetTimeByIdate(idate string) int64 {
	l := len(idate) //不检查 idate
	return TimeStrToTime(idate[0:l-4] + "-" + idate[l-4:l-2] + "-" + idate[l-2:l] + " 00:00:00")
}

//idate 转 0点时间戳
func GetStrTimeByIdate(idate string) string {
	l := len(idate) //不检查 idate
	return (idate[0:l-4] + "-" + idate[l-4:l-2] + "-" + idate[l-2:l] + " 00:00:00")
}

// 当前时间距离今日零点的差值 5分钟
func GetTimeMinutr5Diff() int64 {
	t := time.Now().Unix()
	diff := t - GetTimeBegin(t)
	return diff / 60 / 5
}

// 传入指定时间，返回当前时间的零点
func GetTimeBegin(t int64) int64 {
	str := time.Unix(t, 0).Format("2006-01-02")
	tt, _ := time.ParseInLocation(TimeFormat_Sec, str+" 00:00:00", time.Local)
	return tt.Unix()
}

func Second2days(nSec int64) string {
	mins := int64(0)
	hours := int64(0)
	days := int64(0)

	if nSec >= 60 {
		mins = nSec / 60
		nSec = nSec % 60
	}
	if mins >= 60 {
		hours = mins / 60
		mins = mins % 60
	}
	if hours >= 24 {
		days = hours / 24
		hours = hours % 24
	}
	output := ""
	if days > 0 {
		output = fmt.Sprintf("%dd", days)
	}
	if hours > 0 {
		output += fmt.Sprintf("%dh", hours)
	}
	if mins > 0 {
		output += fmt.Sprintf("%dm", mins)
	}
	return output
}
