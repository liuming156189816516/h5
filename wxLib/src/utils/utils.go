package utils

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
	"io"
	"math/rand"
	"net"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

//进制 字节
const (
	BASE_TEN     = 10 //10进制
	BASE_TWO     = 2  //2进制
	BASE_SIXTEEN = 16

	BIT_THIRTY_TWO = 32 //32位
	BIT_SIXTY_FOUR = 64 //64位
)

var billSeq uint64

// Substr returns the substr from start to length.
func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	if start > len(bt) {
		start = start % len(bt)
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
	} else {
		end = start + length
	}
	return string(bt[start:end])
}
func GetUUID() string {
	u, err := uuid.NewV4()
	if err != nil {
		//logs.LogError("GetUUID err:%+v", err)
		return ""
	}
	return fmt.Sprintf("%s", u)
}

func ReplaceSpecial(s string) string {
	chars := []string{"]", "%", "'", "^", "\\\\", "[", ".", "(", ")", "-"}
	r := strings.Join(chars, "")
	re := regexp.MustCompile("[" + r + "]+")
	return re.ReplaceAllString(s, "")
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

/**
判断手机系统版本
conVer 条件版本号
mVer 当前版本号
*/
func CheckAppVersion(conVer, mVer string) bool {
	conVersion, _ := strconv.ParseInt(conVer, 10, 64)
	// logs.LogDebug("CheckAppVersion conVersion:%d", conVersion)
	if conVersion <= 0 {
		return true
	}
	mVersion, _ := strconv.ParseInt(mVer, 10, 64)
	// logs.LogDebug("CheckAppVersion conVersion:%d,mVersion:%d", conVersion, mVersion)
	return mVersion >= conVersion
}

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func IdInIds(id int32, ids []int32) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}

func StrIdInIds(id string, ids []string) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}

//int 转 string
func IntToStr(data int64) string {
	return strconv.FormatInt(data, BASE_TEN)
}

//uint 转 string
func UintToStr(data uint64) string {
	return strconv.FormatUint(data, BASE_TEN)
}

//string 转 int  //出错 就是 0 用的时候 注意
func StrToInt64(data string) int64 {
	ret := int64(0)
	ret, err := strconv.ParseInt(data, BASE_TEN, BIT_SIXTY_FOUR)
	if err != nil {
		return ret
	}
	return ret
}

func MapToInt64(m map[string]string, key string, def int64) int64 {
	data, _ := m[key]
	if data == "" {
		return def
	}
	n := StrToInt64(data)
	return n
}

func StrToInt(src string) int {
	_val, _ := strconv.ParseInt(src, BASE_TEN, BIT_SIXTY_FOUR)
	return int(_val)
}

func StrMapToInt(m map[string]string, key string) int {
	if m == nil {
		return 0
	}
	key = strings.ToLower(key)
	src, ok := m[key]
	if !ok {
		return 0
	}
	_val, _ := strconv.ParseInt(src, BASE_TEN, BIT_SIXTY_FOUR)
	return int(_val)
}

//string 转 int  //出错 就是 0 用的时候 注意
func StrToUInt(data string) uint64 {
	ret := int64(0)
	ret, err := strconv.ParseInt(data, BASE_TEN, BIT_SIXTY_FOUR)
	if err != nil {
		return uint64(ret)
	}
	return uint64(ret)
}

//16进制string 转 int  //出错 就是 0 用的时候 注意
func Str16ToInt(data string) int64 {
	ret := int64(0)
	ret, err := strconv.ParseInt(data, BASE_SIXTEEN, BIT_SIXTY_FOUR)
	if err != nil {
		return ret
	}
	return ret
}

func InArrayInt(arr []int, x int) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

func InsertArrayInt64(inArr []int64, index int, value int64) []int64 {
	rear := append([]int64{}, inArr[index:]...)
	inArr = append(inArr[0:index], value)
	inArr = append(inArr, rear...)
	return inArr
}
func InArrayInt64(arr []int64, x int64) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

func InsertArrayUint64(inArr []uint64, index int, value uint64) []uint64 {
	rear := append([]uint64{}, inArr[index:]...)
	inArr = append(inArr[0:index], value)
	inArr = append(inArr, rear...)
	return inArr
}
func InArrayUint64(arr []uint64, x uint64) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

func StrToFloat(data string) float64 {
	ret := float64(0)
	ret, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return ret
	}
	return ret
}

func InsertArrayInt32(inArr []int32, index int, value int32) []int32 {
	rear := append([]int32{}, inArr[index:]...)
	inArr = append(inArr[0:index], value)
	inArr = append(inArr, rear...)
	return inArr
}
func InArrayInt32(arr []int32, x int32) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

func InArrayString(arr []string, x string) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}
func StringHasPrefix(x string, arr []string) bool {
	for _, v := range arr {
		if strings.HasPrefix(x, v) {
			return true
		}
	}
	return false
}

func StringHasSub(x string, arr []string) bool {
	for _, v := range arr {
		if strings.Index(x, v) >= 0 {
			return true
		}
	}
	return false
}

func GetRandString(arr []string) string {
	if len(arr) <= 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(arr))
	return arr[i]
}

// min 到 max之间的随机数
func RandGen(min, max int) int {
	if max <= min { //最大小于最小
		return min
	}
	max = max - min + 1
	return rand.Intn(max) + min
}
func RandGenint64(min, max int64) int64 {
	if max <= min { //最大小于最小
		return min
	}
	max = max - min + 1
	return rand.Int63n(max) + min
}

func SplitToInt(s string, sep string) []int {
	arrStr := strings.Split(s, sep)
	arr := make([]int, 0, len(arrStr))
	for _, v := range arrStr {
		if v == "" {
			continue
		}
		x, _ := strconv.Atoi(v)
		if x > 0 {
			arr = append(arr, x)
		}
	}
	return arr
}

func SplitToInt32(s string, sep string) []int32 {
	arrStr := strings.Split(s, sep)
	arr := make([]int32, 0, len(arrStr))
	for _, v := range arrStr {
		if v == "" {
			continue
		}
		x, _ := strconv.Atoi(v)
		if x > 0 {
			arr = append(arr, int32(x))
		}
	}
	return arr
}

// 这个函数慎用 obj 不能用指针
func Struct2Map(obj interface{}) map[string]interface{} {
	var data = make(map[string]interface{})
	if obj == nil {
		return data
	}
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	for i := 0; i < t.NumField(); i++ {
		data[ /*strings.ToLower*/ (t.Field(i).Name)] = v.Field(i).Interface()
	}
	return data
}

//用map填充结构
func FillStruct(data map[string]interface{}, obj interface{}) error {
	for k, v := range data {
		err := SetField(obj, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

//用map的值替换结构的值
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()        //结构体属性值
	structFieldValue := structValue.FieldByName(name) //结构体单个属性值

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type() //结构体的类型
	val := reflect.ValueOf(value)              //map值的反射值

	var err error
	if structFieldType != val.Type() {
		val, err = TypeConversion(fmt.Sprintf("%v", value), structFieldValue.Type().Name()) //类型转换
		if err != nil {
			return err
		}
	}

	structFieldValue.Set(val)
	return nil
}

func AutoToString(a interface{}) string {
	if a == nil {
		return ""
	}
	switch s := a.(type) {
	case string:
		return (s)
	case *string:
		return (*s)
	case bool:
		return strconv.FormatBool(s) //fmt.Sprintf("%t", s)
	case *bool:
		return strconv.FormatBool(*s) //fmt.Sprintf("%t", *s)
	case int:
		return strconv.FormatInt(int64(s), 10)
	case *int32:
		return strconv.FormatInt(int64(*s), 10) //fmt.Sprintf("%d", *s)
	case *int64:
		return strconv.FormatInt((*s), 10) //fmt.Sprintf("%d", *s)
	case float64:
		return fmt.Sprintf("%g", s)
	case *float64:
		return fmt.Sprintf("%g", *s)
	default:
		return ""
	}
}

//类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//else if .......增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}

//按百分比检查概率
func CheckRateByPer(per int) bool {
	if per <= 0 {
		return false
	}
	if per >= 100 {
		return true
	}
	n := rand.Intn(100)
	return n < per
}

//按分母检查概率
func CheckRateByBase(b int) bool {
	if b < 1 {
		return false
	}
	return rand.Intn(b) == 0
}

func IdsIntToStr(ids []int32) string {
	str := ""
	for _, v := range ids {
		if str == "" {
			str += fmt.Sprintf("%d", v)
		} else {
			str += fmt.Sprintf(",%d", v)
		}
	}
	return str
}

func IdsStrToStr(ids []string) string {
	str := ""
	for _, v := range ids {
		if str == "" {
			str += fmt.Sprintf(" '%v' ", v)
		} else {
			str += fmt.Sprintf(" ,'%v' ", v)
		}
	}
	return str
}

func StrInStrs(str string, strs []string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

func GenBillNo(prefixFmt string, args ...interface{}) string {
	i := atomic.AddUint64(&billSeq, 1)
	prefix := fmt.Sprintf(prefixFmt, args...)

	return fmt.Sprintf("%s|%03d", prefix, i)
}

func IsSlic(data interface{}) bool {
	if data == nil {
		return false
	}
	v := reflect.ValueOf(data)
	k := v.Kind()
	if k == reflect.Slice {
		return true
	}
	return false
}

//获取字符串 如果不是字符串 返回""
func GetString(data interface{}) string {
	if data == nil {
		return ""
	}
	v := reflect.ValueOf(data)
	k := v.Kind()
	if k == reflect.String {
		if t, ok := data.(string); ok {
			return t
		} else { //bson 字符串
			if t, ok := data.(bson.ObjectId); ok {
				return t.Hex()
			}
		}
	}
	if k == reflect.Ptr {
		t, ok := data.(*string)
		if ok {
			return *t
		}
	}
	return ""
}

//获取int64 如果不是numbel 类型 返回0
func GetInt64(data interface{}) int64 {
	if data == nil {
		return 0
	}
	v := reflect.ValueOf(data)
	k := v.Kind()
	switch k {
	case reflect.Int:
		if t, ok := data.(int); ok {
			return int64(t)
		}
	case reflect.Int8:
		if t, ok := data.(int8); ok {
			return int64(t)
		}
	case reflect.Int16:
		if t, ok := data.(int16); ok {
			return int64(t)
		}
	case reflect.Int32:
		if t, ok := data.(int32); ok {
			return int64(t)
		}
	case reflect.Int64:
		if t, ok := data.(int64); ok {
			return int64(t)
		}
	case reflect.Uint:
		if t, ok := data.(uint); ok {
			return int64(t)
		}
	case reflect.Uint8:
		if t, ok := data.(uint8); ok {
			return int64(t)
		}
	case reflect.Uint16:
		if t, ok := data.(uint16); ok {
			return int64(t)
		}
	case reflect.Uint32:
		if t, ok := data.(uint32); ok {
			return int64(t)
		}
	case reflect.Uint64:
		if t, ok := data.(uint64); ok {
			return int64(t)
		}
	case reflect.Ptr:
		if t, ok := data.(*int); ok {
			return int64(*t)
		}
		if t, ok := data.(*int8); ok {
			return int64(*t)
		}
		if t, ok := data.(*int16); ok {
			return int64(*t)
		}
		if t, ok := data.(*int32); ok {
			return int64(*t)
		}
		if t, ok := data.(*int64); ok {
			return int64(*t)
		}
		if t, ok := data.(*uint); ok {
			return int64(*t)
		}
		if t, ok := data.(*uint8); ok {
			return int64(*t)
		}
		if t, ok := data.(*uint16); ok {
			return int64(*t)
		}
		if t, ok := data.(*uint32); ok {
			return int64(*t)
		}
		if t, ok := data.(*uint64); ok {
			return int64(*t)
		}
	}
	return 0
}
//获取int64 如果不是numbel 类型 返回0
func GetFloat64(data interface{}) float64 {
	if data == nil {
		return 0
	}
	v := reflect.ValueOf(data)
	k := v.Kind()
	switch k {
	case reflect.Float64:
		if t, ok := data.(float64); ok {
			return float64(t)
		}
	}
	return 0
}


//读取csv 文件
func ReadCsvData(data string, title ...string) map[string][]string {
	rfile := strings.NewReader(data)
	selfcsv := csv.NewReader(rfile)
	ret := map[string][]string{}
	if selfcsv == nil || len(title) == 0 {
		return ret
	}
	i := int32(0)
	tm := map[string]int{}
	for _, t := range title {
		tm[t] = 0
		ret[t] = []string{}
	}
	idm := map[int]string{}
	for {
		tmp, err := selfcsv.Read()
		if err != nil {
			if err != io.EOF {

			}
			break
		}
		if i == 0 { //是title
			for id, ti := range tmp {
				if _, ok := tm[ti]; ok {
					tm[ti] = id //这个标题在第几行
					idm[id] = ti
				}
			}
		} else {
			for id, dd := range tmp {
				if ti, ok := idm[id]; ok { //有需要
					if old, ok := ret[ti]; ok {
						old = append(old, dd)
						ret[ti] = old
					} else {
						ret[ti] = []string{dd}
					}
				}
			}
		}
		i++
	}
	return ret
}

//是否是本机Ip
func IsLocalIP(ip string) (bool, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false, err
	}
	for i := range addrs {
		intf, _, err := net.ParseCIDR(addrs[i].String())
		if err != nil {
			return false, err
		}
		if net.ParseIP(ip).Equal(intf) {
			return true, nil
		}
	}
	return false, nil
}

func TrimStr(str string) string {
	str = strings.Trim(str, " ")
	str = strings.Trim(str, "\n")
	str = strings.Trim(str, "\r")
	return str
}

func GenUrl(host string, port int32, usPort int32) string {
	// if strings.HasPrefix(host, "http") == false || strings.HasPrefix(host, "HTTP") == false {
	// 	host = "https://" + host
	// 	if !(port == 443 || port == 80) {
	// 		host += ":" + strconv.Itoa(int(port))
	// 	}
	// }
	host = strings.TrimSpace(host)
	if port > 0 {
		return fmt.Sprintf("https://%s:%d", host, port)
	}
	if usPort > 0 {
		return fmt.Sprintf("http://%s:%d", host, usPort)
	}

	return ""
}

//  10005 / 10000  = 1 返回 true
func CheckIntMold(a int64, str string) bool {
	b, _ := strconv.ParseInt(str, 10, 64)
	if a <= 0 || b <= 0 {
		return false
	}
	return b/a == 1
}

func GetOsBuildId(str string) string {
	if str == "" {
		return ""
	}
	s := strings.Split(str, "|")
	if len(s) > 2 {
		return s[2]
	}
	return ""
}

// 数组去重
func GetArrayUnique(data []int32) []int32 {
	out := []int32{}
	for _, v := range data {
		if IdInIds(v, out) == false {
			out = append(out, v)
		}
	}
	return out
}

func StrInArray(str string, array []string) bool {
	if str == "" || len(array) <= 0 {
		return false
	}
	for _, v := range array {
		if v != "" && str == v {
			return true
		}
	}
	return false
}

func UniqueIds(ids []int32) []int32 {
	if len(ids) <= 0 {
		return ids
	}
	out := []int32{}
	for _, v := range ids {
		if !IdInIds(v, out) {
			out = append(out, v)
		}
	}
	return out
}

func UniqueIdsInt64(ids []int64) []int64 {
	if len(ids) <= 0 {
		return ids
	}
	out := []int64{}
	for _, v := range ids {
		if !IdInIdsINT64(v, out) {
			out = append(out, v)
		}
	}
	return out
}

func IdInIdsINT64(id int64, ids []int64) bool {
	for _, v := range ids {
		if id == v {
			return true
		}
	}
	return false
}

func UniqueStrIds(ids []string) []string {
	if len(ids) <= 0 {
		return ids
	}
	out := []string{}
	for _, v := range ids {
		if !StrIdInIds(v, out) {
			out = append(out, v)
		}
	}
	return out
}

func RemoveArrayItem(arr []int32, remove []int32) []int32 {
	if len(remove) <= 0 {
		return arr
	}
	data := []int32{}
	for _, v := range arr {
		for _, k := range remove {
			if v == k {
				continue
			}
		}
		data = append(data, v)
	}
	return data
}

//切片乱序
func RandSlice(slice interface{}) {
	rv := reflect.ValueOf(slice)
	if rv.Type().Kind() != reflect.Slice {
		return
	}

	length := rv.Len()
	if length < 2 {
		return
	}

	swap := reflect.Swapper(slice)
	rand.Seed(time.Now().Unix())
	for i := length - 1; i >= 0; i-- {
		j := rand.Intn(length)
		swap(i, j)
	}
	return
}

//字符首字母大写转换
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

//字符串数组乱序
func RandomStrList(strings []string) []string { //字符串数组
	for i := len(strings) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		strings[i], strings[num] = strings[num], strings[i]
	}
	str2 := []string{}
	for i := 0; i < len(strings); i++ {
		str2 = append(str2, strings[i])
	}
	return str2
}

//get英文名字
func GetEFullName() string {
	var lastName = []string{"Brenda", "Young", "Carol", "Davis", "Sharon", "Harris", "Scott", "Clark", "Sharon", "Garcia", "Kimberly", "Davis", "William", "Harris", "Eric", "Martinez", "Melissa", "Martin", "Ronald", "Thomas", "Michelle", "Harris", "Charles", "Clark", "Charles", "Miller", "Mark", "Wilson", "Amy", "White", "Jeffrey", "Miller", "Sharon", "Taylor", "Charles", "Anderson", "Jason", "Hall", "Helen", "Lee", "William", "Miller", "Linda", "Garcia", "Kevin", "Allen", "Deborah", "Gonzalez", "Sandra", "Miller", "Angela", "Thompson", "Scott", "Davis", "Jessica", "Perez", "Laura", "Lopez", "Anna", "Johnson", "Sharon", "White", "Christopher", "Martin", "Angela", "Lopez", "David", "Smith", "Paul", "Gonzalez", "Daniel", "Thompson", "Scott", "Williams", "Richard", "Johnson", "Thomas", "Gonzalez", "Margaret", "Miller", "Jose", "Walker", "Sarah", "Lee", "Thomas", "Garcia", "Mark", "Martinez", "David", "Anderson", "Eric", "Thomas", "Jason", "Davis", "Brian", "Smith", "Sarah", "Davis", "John", "Lee", "Richard", "Smith", "James", "Jones", "Gary", "Johnson", "Daniel", "Jackson", "Susan", "Smith", "Melissa", "Hall", "Daniel", "Garcia", "Patricia", "Rodriguez", "Frank", "Allen", "Betty", "Young", "Linda", "Perez", "Paul", "Clark", "Lisa", "Perez", "Karen", "Thompson", "Michael", "Garcia", "Betty", "White", "Angela", "Rodriguez", "John", "Hall", "Angela", "Perez", "Amy"}
	var firstName = []string{"Walker", "Kimberly", "Williams", "James", "Rodriguez", "David", "Gonzalez", "Robert", "Jackson", "Steven", "Brown", "Ronald", "Smith", "Carol", "Rodriguez", "Barbara", "Thompson", "Charles", "Davis", "Jennifer", "Taylor", "Matthew", "Rodriguez", "Laura", "Brown", "Paul", "Brown", "Barbara", "Wilson", "Jessica", "Jackson", "Karen", "Clark", "David", "Williams", "Dorothy", "Walker", "Shirley", "Clark", "Linda", "Young", "Nancy", "Thomas", "Carol", "Hall", "David", "Rodriguez", "Barbara", "Hall", "Gary", "Rodriguez", "Anna", "Anderson", "Mark", "Lopez", "Frank", "Lewis", "Helen", "Taylor", "Amy", "Hall", "Sandra", "Johnson", "Kenneth", "Clark", "Matthew", "Perez", "Richard", "Williams", "Nancy", "Hernandez", "Larry", "Miller", "Donna", "Martin", "Karen", "Lopez", "Donald", "Clark", "Edward", "White", "Matthew", "Clark", "William", "Smith", "Donna", "Thompson", "Brenda", "Taylor", "Lisa", "Williams", "Frank", "Davis", "Deborah", "White", "Steven", "Lewis", "Daniel", "Thomas", "Anthony", "Robinson", "Ruth", "White", "Thomas", "Robinson", "David", "Thomas", "Ruth", "Williams", "Lisa", "Hernandez", "Barbara", "Clark", "Sharon", "Johnson", "Margaret", "Wilson", "Gary", "Young", "Linda", "Gonzalez", "Carol", "Harris", "Brian", "Hall", "Amy", "Gonzalez", "Matthew", "Lopez", "Sandra", "Walker", "Charles", "Lee", "Angela", "Johnson", "Joseph", "Wilson", "Edward", "Hall", "Paul"}
	var lastNameLen = len(lastName)
	var firstNameLen = len(firstName)
	rand.Seed(time.Now().UnixNano())     //设置随机数种子
	var first string                     //名
	for i := 0; i <= rand.Intn(1); i++ { //随机产生2位或者3位的名
		first = fmt.Sprint(firstName[rand.Intn(firstNameLen-1)])
	}
	//返回姓名
	return fmt.Sprintf("%s%s", fmt.Sprint(lastName[rand.Intn(lastNameLen-1)]), first)
}

//get中文名字
func GetZFullName() string {
	var lastName = []string{
		"赵", "钱", "孙", "李", "周", "吴", "郑", "王", "冯", "陈", "褚", "卫", "蒋",
		"沈", "韩", "杨", "朱", "秦", "尤", "许", "何", "吕", "施", "张", "孔", "曹", "严", "华", "金", "魏",
		"陶", "姜", "戚", "谢", "邹", "喻", "柏", "水", "窦", "章", "云", "苏", "潘", "葛", "奚", "范", "彭",
		"郎", "鲁", "韦", "昌", "马", "苗", "凤", "花", "方", "任", "袁", "柳", "鲍", "史", "唐", "费", "薛",
		"雷", "贺", "倪", "汤", "滕", "殷", "罗", "毕", "郝", "安", "常", "傅", "卞", "齐", "元", "顾", "孟",
		"平", "黄", "穆", "萧", "尹", "姚", "邵", "湛", "汪", "祁", "毛", "狄", "米", "伏", "成", "戴", "谈",
		"宋", "茅", "庞", "熊", "纪", "舒", "屈", "项", "祝", "董", "梁", "杜", "阮", "蓝", "闵", "季", "贾",
		"路", "娄", "江", "童", "颜", "郭", "梅", "盛", "林", "钟", "徐", "邱", "骆", "高", "夏", "蔡", "田",
		"樊", "胡", "凌", "霍", "虞", "万", "支", "柯", "管", "卢", "莫", "柯", "房", "裘", "缪", "解", "应",
		"宗", "丁", "宣", "邓", "单", "杭", "洪", "包", "诸", "左", "石", "崔", "吉", "龚", "程", "嵇", "邢",
		"裴", "陆", "荣", "翁", "荀", "于", "惠", "甄", "曲", "封", "储", "仲", "伊", "宁", "仇", "甘", "武",
		"符", "刘", "景", "詹", "龙", "叶", "幸", "司", "黎", "溥", "印", "怀", "蒲", "邰", "从", "索", "赖",
		"卓", "屠", "池", "乔", "胥", "闻", "莘", "党", "翟", "谭", "贡", "劳", "逄", "姬", "申", "扶", "堵",
		"冉", "宰", "雍", "桑", "寿", "通", "燕", "浦", "尚", "农", "温", "别", "庄", "晏", "柴", "瞿", "阎",
		"连", "习", "容", "向", "古", "易", "廖", "庾", "终", "步", "都", "耿", "满", "弘", "匡", "国", "文",
		"寇", "广", "禄", "阙", "东", "欧", "利", "师", "巩", "聂", "关", "荆", "司马", "上官", "欧阳", "夏侯",
		"诸葛", "闻人", "东方", "赫连", "皇甫", "尉迟", "公羊", "澹台", "公冶", "宗政", "濮阳", "淳于", "单于",
		"太叔", "申屠", "公孙", "仲孙", "轩辕", "令狐", "徐离", "宇文", "长孙", "慕容", "司徒", "司空"}
	var firstName = []string{
		"伟", "刚", "勇", "毅", "俊", "峰", "强", "军", "平", "保", "东", "文", "辉", "力", "明", "永", "健", "世", "广", "志", "义",
		"兴", "良", "海", "山", "仁", "波", "宁", "贵", "福", "生", "龙", "元", "全", "国", "胜", "学", "祥", "才", "发", "武", "新",
		"利", "清", "飞", "彬", "富", "顺", "信", "子", "杰", "涛", "昌", "成", "康", "星", "光", "天", "达", "安", "岩", "中", "茂",
		"进", "林", "有", "坚", "和", "彪", "博", "诚", "先", "敬", "震", "振", "壮", "会", "思", "群", "豪", "心", "邦", "承", "乐",
		"绍", "功", "松", "善", "厚", "庆", "磊", "民", "友", "裕", "河", "哲", "江", "超", "浩", "亮", "政", "谦", "亨", "奇", "固",
		"之", "轮", "翰", "朗", "伯", "宏", "言", "若", "鸣", "朋", "斌", "梁", "栋", "维", "启", "克", "伦", "翔", "旭", "鹏", "泽",
		"晨", "辰", "士", "以", "建", "家", "致", "树", "炎", "德", "行", "时", "泰", "盛", "雄", "琛", "钧", "冠", "策", "腾", "楠",
		"榕", "风", "航", "弘", "秀", "娟", "英", "华", "慧", "巧", "美", "娜", "静", "淑", "惠", "珠", "翠", "雅", "芝", "玉", "萍",
		"红", "娥", "玲", "芬", "芳", "燕", "彩", "春", "菊", "兰", "凤", "洁", "梅", "琳", "素", "云", "莲", "真", "环", "雪", "荣",
		"爱", "妹", "霞", "香", "月", "莺", "媛", "艳", "瑞", "凡", "佳", "嘉", "琼", "勤", "珍", "贞", "莉", "桂", "娣", "叶", "璧",
		"璐", "娅", "琦", "晶", "妍", "茜", "秋", "珊", "莎", "锦", "黛", "青", "倩", "婷", "姣", "婉", "娴", "瑾", "颖", "露", "瑶",
		"怡", "婵", "雁", "蓓", "纨", "仪", "荷", "丹", "蓉", "眉", "君", "琴", "蕊", "薇", "菁", "梦", "岚", "苑", "婕", "馨", "瑗",
		"琰", "韵", "融", "园", "艺", "咏", "卿", "聪", "澜", "纯", "毓", "悦", "昭", "冰", "爽", "琬", "茗", "羽", "希", "欣", "飘",
		"育", "滢", "馥", "筠", "柔", "竹", "霭", "凝", "晓", "欢", "霄", "枫", "芸", "菲", "寒", "伊", "亚", "宜", "可", "姬", "舒",
		"影", "荔", "枝", "丽", "阳", "妮", "宝", "贝", "初", "程", "梵", "罡", "恒", "鸿", "桦", "骅", "剑", "娇", "纪", "宽", "苛",
		"灵", "玛", "媚", "琪", "晴", "容", "睿", "烁", "堂", "唯", "威", "韦", "雯", "苇", "萱", "阅", "彦", "宇", "雨", "洋", "忠",
		"宗", "曼", "紫", "逸", "贤", "蝶", "菡", "绿", "蓝", "儿", "翠", "烟", "小", "轩"}
	var lastNameLen = len(lastName)
	var firstNameLen = len(firstName)
	rand.Seed(time.Now().UnixNano())     //设置随机数种子
	var first string                     //名
	for i := 0; i <= rand.Intn(1); i++ { //随机产生2位或者3位的名
		first = fmt.Sprint(firstName[rand.Intn(firstNameLen-1)])
	}
	//返回姓名
	return fmt.Sprintf("%s%s", fmt.Sprint(lastName[rand.Intn(lastNameLen-1)]), first)
}

//获取随机字符串
func RandStr(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}


//获取随机字符串
func RandNum(n int) string {
	var letters = []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//替换字符串前缀
func ReplaceStringPrefix(original, oldPrefix string) string {
	if strings.HasPrefix(original, oldPrefix) {
		return  original[len(oldPrefix):]
	}
	return original
}

//获取随机表情
func RandStr2(n int) string {
	var letters = []rune("🤱🙇‍♀️🚴🏽‍♀️🤙👮‍♂️👨💄👨‍👩‍👧🦿🍌💶🧘🏾‍♂️🪦🏌️‍♀️🌭⛱👳‍♀️💆‍♂️🤭🙈🐏🍾🔮🛍👨‍⚕️⚒😉👿🐼🧄🎍📔📪👨‍👨‍👧‍👦🏋️‍♂️🚴🏽‍♂️🖤🕺🍣🔪🚌👔💇‍♂️😐🍖👨‍👨‍👧😹🦮🌕🌌⛈🐞⏳👞⛳️💛🤌🏄🏎🚿👨‍🔧👩‍🔬🚣🏿‍♂️💯🐍🕌⚾🗳🧾🍋🍕🕵️‍♀️💰💙🗨🐖🥪☕🕤🎽🤦‍♀️🌚♦🎥🤾🏾‍♀️🤖💁🌃🎸📬🚴🏻‍♀️👁🧖👩‍🏭🏄🏽‍♀️😸🦖👛🔉👨‍👧‍👦☘️🤾🏽‍♂️😏💅🧵💂‍♀️👆🧒👸🐐🎮⌚️🙇‍♂️💕🙆🦹🍏🏡🕙🌪🧗🏼‍♂️😂😫🐃🍑🌄🔔📖🫀🐊🥐🏋🏾‍♂️🕒😃🤣💢💤🤛🍷⏰🔌🪝😝🍡🧜‍♂️🙆‍♂️🤸🏽‍♂️🏌🏽‍♂️🙅🚣🥦💎🧘🏽‍♀️🤾🏾‍♂️🫓🎪🚊🚞🔒✍️🧚‍♀️🧗🏾‍♀️🥸🐕🚇🎲🖲📨⚱🧜🕗🪖🧞🧉🚂🌨🎈🏈🧦🙏⛷🛢📜👨‍👦‍👦⚾️🚵🏿‍♀️😘👨‍✈️🧗🏾‍♂️🥰🤞👍🚗🚪🏌️‍♂️🚣‍♂️🧐🙊👼🚴🐄🥧⚖💋🌯👩‍👦‍👦⛹🏼‍♂️👤🎵🎛🗣🍧🗂☺️👨‍🏫🤽🏼‍♀️🏌🏼‍♀️💖👶🧏🎳🧩🔖👨‍🚀🤽‍♂️🌋🛷💂‍♂️👨‍🍳🧘🏽‍♂️🦙🥅🎓🪕⛏✌🧙🐗🛬🖌🗡🥵📆📉☝🕰🪁📁🪰🪱🌔😚🛀🎣🧴🤹🏿‍♀️🤑😟🚋🌈☎️💪🛌🐒📭👩‍❤️‍👩🏄🏼‍♀️⌨️🦧🐷⛹🏿‍♀️🍍🏫🪃🪤🏇🏻🤥👂🚔💒🥲😔😕😤☠🍈🏘📏⚗☕️🏄‍♀️🚴‍♂️🤟👵🐀🍸🕧👚🌎💧🏷💪🖕👲🍀🍥🌗🥎📘💼👨‍🔬⚡️🤽🏽‍♀️🥼📠📮🪛🤼🩸👩‍🎤🕠🖊🥷🏸🤸🏼‍♀️🤾🏻‍♀️🏌🏽‍♀️🚣🏼‍♂️🤪🕋🕢🎱📕🤧👯‍♀️🪶🌓🦺📫🔏👨‍👧⛹🏻‍♂️💑🍃🌇🥈🔧🏄🏿‍♀️🌼🍎🥓🏏🎯🗞🚴🏼‍♂️👯🧗🥨🚡🛒😞🤬🍔⛽🛶🥾🏋🏾‍♀️🤡🧠👪🍼🍴🪑🤸🏼‍♂️🥂🗒🩺😼🧚🫒🥫🌘🧿🧥🍒🎖😊🏔📢🧬☀️👨‍💻👁‍🗨🐽🌷🍤🚏🚨👒🛹🏋🏼‍♂️⛹🏾‍♂️🏋🙋‍♀️🍩🕘🎻🚴🏿‍♂️😩🍟🏩🔗🗿🍶⌨🔑🤸‍♀️🏊🏿‍♂️🧛💮🎴😆🙋🦢🫖🏕🥿💀🦂🍗🍜🌆🤹🏾‍♂️😈👐🎞🗜🧷🏄🏼‍♂️🤹🏽‍♀️👨‍⚖️😶🐬🥭🌏🚦👨‍👨‍👦‍👦🤾🏼‍♀️🐨🦞🏉👓🤦🕷🌐📷👨‍🎓👩‍💻🧣😒🕳🙎🍚🦑🗽🏓🪞⛹🏾‍♀️☹🦝🧧👨‍❤️‍💋‍👨🏊🏼‍♀️🍮🏥🚓♠🧞‍♀️💆‍♀️⭐️😑🏌🕹🔍🧮🧝‍♀️👩‍👩‍👦🤷👭👢🧙‍♀️🍇🛡🧘‍♂️👄🧱🛖🛳🏐🥁🤹‍♀️🧆🪡🛋🧺🏄🏽‍♂️🚣🏽‍♀️🚣🏾‍♂️💘🤵🐉🪢📐🤹🏾‍♀️🤽🦃🍠🪀🧪🤸🏾‍♂️🤫📟🤽🏼‍♂️🧎📣🪓🚲✂🧗🏻‍♂️🌍🛼🛣🕖🧝🥟🏝🌝🌞📃🦆🛥🗄✈️🐆🐇🌽🚘🖋👩‍👧‍👦🚵‍♀️😥🐪⛵🌒🎋👨‍🎤🧛‍♂️📻😍🍯🌛🌠☄🎎🩳👩‍🎨🦥🍆🩲📒👩‍👩‍👦‍👦👨‍👧‍👧🚵🏻‍♀️😖🦵🔐🕵️‍♂️🤾🏻‍♂️🚴🏿‍♀️🤾‍♀️🤜🦗🥩🏬⛑🎧💊👉🕊🌳🌖🌦👩‍👩‍👧‍👧🍪🧯💁‍♀️😠👊🧋🏣👨‍👩‍👧‍👦🏋🏻‍♂️🧍🦁🐙🫐🏖⚓️🥺*‍❄📍👨‍🎨🤓🐓🌶🕛🚴‍♀️🦏🪅🔈💾🔩🙍‍♂️🏇🏼💈🚥👩‍⚖️🧜‍♀️🏋🏻‍♀️🧗🏽‍♀️🤾🏼‍♂️🏠🏤⛅👗🎩📝🏋🏽‍♀️😽❤🤎🖖🍨💹🦛🍂🥖💺🕕🌧🥌🗝🧼🚵🏽‍♀️😳🌲🥞🍵🕍🤿🏺☃🥋🪆😿💥👎🧀☄️🚴🏼‍♀️👏👮🐝🍰🍽⛲🚃🐦🐋👖🪥🏇🏾⚱️😭👣⛴👨‍🚒👩‍🏫😛👻💬🐤🫕🧳🎟🚴🏾‍♀️💏🦤🛤⛹️‍♀️🧗‍♀️⌛️🏊‍♂️😜💫🤏💃🍅🚎🛵🤽🏻‍♂️🚵🏿‍♂️🤹‍♂️😄🗾🦼😮🐸🏢🤩🛫📑⛲️🗻⛸🦯🏃‍♀️👨‍👦⛹🏿‍♂️🧘🏻‍♀️😎🐈‍⬛🕥⚙🎅🐟🍞🚉👕🪟🧝‍♂️😵🧡🦍🗺🚑🛸☝️🏊🏽‍♀️🏊🏾‍♀️🤹🏻‍♂️👽🌱🥕🛎⏲🤽🏽‍♂️🤤🤠🐣🚚🚤📽🗯🥀💴🧙‍♂️⛪️🤮☂📧🛁💓🐲⛩☁️🪨☀🎰🏋️‍♀️🏌🏼‍♂️🤴🍙🕚⚡🎺🔦⛹🏼‍♀️💔⛹🍢🪒⚙️💝👇🦅🦽🎇💵🎤😯💇🐿🧁🏦👘🩴👩‍👦🏄🏾‍♀️🧚‍♂️🐅🧭🌙🎃🏆👜📦😢🎚🚶‍♂️🤸‍♂️👴🐘🚄🏍🏃‍♂️💭🧑🦣🌺🏋🏿‍♀️🤒👹👈⏱📩👨‍👨‍👧‍👧🤹🏼‍♀️😡🚅🔓🔬🙀💗🦾🐩🧃🏭🌁🏃🐹🥠☁📎🚣🏿‍♀️🚵🏻‍♂️🏊🏻‍♀️🥳🤺🧸👑👷‍♂️👩‍🚒🤸🏿‍♀️🤽🏻‍♀️🚣🏻‍♀️🧗🏻‍♀️💆🥏🕯🤾‍♂️🧗‍♂️😬🦷🚵🥚🩰🦩💻💳🖍⛹️‍♂️🏄‍♂️🤹🐾🌜💇‍♀️🧗🏽‍♂️🚵‍♂️🙄🤯🫁🤸🦈🛻📋🤘🪴🍲🍘🛠🥶🕑🧨🧲🙋‍♂️🪂🏹💁‍♂️🤹🏼‍♂️😻😾🌥🥉🤽🏿‍♀️😁🦶🦟🕣📥🧫👧👱🥥🥄🚟🪐📙🏄🏻‍♂️✌️😌🧘🐈🏞🚠🌬🖼🧕🏂🍱🀄🪜🩹🌿🍝🎫🧞‍♂️🤲🤰🐰🛑📱📗🔭👨‍👨‍👦🚵🏼‍♀️⚔️🚛🌟📿📂🗃🧛‍♀️🐺🦬♥👩‍👧‍👧🐁🐳🫔🌅⛳💽👩‍👩‍👧🤔❣🏊🏿‍♀️🤽‍♀️😇🤚🧓🤼‍♀️🙃💂🎭👠🪙🚣🏽‍♂️🦸🕴🦌📄🎡🎐🤐🥢🏪🌑✏⛽️⚰️🐔🫑🖥📼⛹🏽‍♀️📺🏊👥🐭🐥🏯🛕🩱🙅‍♀️🤾🏿‍♂️🚴🏻‍♂️🚆🔊📀🧘🏼‍♂️😺🙍🦨🎗🥇📲👷‍♀️👳🦎💐🚒🕞🎶🐯🐡🚈👝🔕🚣🏾‍♀️🧗🏿‍♀️🍄🥙🥮🏛♨🪗📇🙅‍♂️🤗😦😨😣🕟🎿📞🏇🏽🙌🦦🕶👡🫂🧅👨‍💼💦🏙💞🦀🦪🏋🏼‍♀️⛵️🥱🕵🏰👯‍♂️😅♣🥽🤸🏿‍♂️🧘🏾‍♀️🧽🥴🦴👦👩🦋🌴🧊🧟‍♀️👺👬🦓🦡🪘🌵🧂🍦🪄🎒🎼🖱🏌🏿‍♂️💜💣🥯🥗🛺🎷✒🙉🧇🌩🎾🃏🎹🏇🏿🚶🌸🥔🎠🚖🥍🧘🏿‍♀️😗💟✍👷🌤🌫😱💚🐫🏌🏾‍♀️🚵🏾‍♂️🤕👩‍🌾👩‍🍳🦠🥡🚙🕦⚔📡🚬⚓🙎‍♂️💨🌹🚍🪣🏊‍♀️🤨👩‍👧💌🦕🪲🔇🧻⚰🪧😀👾🦇🍊👱‍♂️🤸🏾‍♀️✊🧔🐱🥘❄🗑🏋🏽‍♂️🥜🕐👙🏌🏻‍♀️🍿🎑💸🙍‍♀️🤸🏽‍♀️🧘🏿‍♂️👅🌻🏗🛩🤦‍♂️🏌🏾‍♂️🚣‍♀️🌮🦐📹👰🦫🦉🦚🕜🪔🧈🥃🖇🧘‍♀️🤽🏾‍♂️😰🦔🏅🤷‍♀️👀🛰🛏🏋🏿‍♂️😷🖐🥛🎆🏮🎂🍻🚐🌡🧤👳‍♂️🍓👩‍⚕️👩‍🔧📰☺🐕‍🦺🍁🕓⛄🌊🔋👩‍🚀🧖‍♂️🤽🏿‍♂️🤹🏿‍♂️⛺️📊🥣👨‍❤️‍👨⛹🏽‍♂️🧗🏼‍♀️🏇👫🥝🧰🤾🏿‍♀️🤹🏽‍♂️🦄🎁🏑✉📈🏊🏽‍♂️🚵🏾‍♀️🖕🤳📚🤾🥑🥤⛰⛅️🚽👩‍👩‍👧‍👦🚵🏼‍♂️💍⛓🏌🏿‍♀️🏵🍛🍭🚴🏾‍♂️🧘🏼‍♀️🦜🥬⛺🕡🏒💿📤♟👨‍🌾🏊🏾‍♂️😋😪🍉🏊🏼‍♂️🚵🏽‍♂️🙂🐵🎢🛴🧶🗓🤸🏻‍♂️😙🐑🏟🥊📯🧖‍♀️👩‍❤️‍💋‍👩👩‍💼🦊🐚🏜⌚🕝🌂🎀⛹🏻‍♀️🪵🏨👱‍♀️👮‍♀️🙎‍♀️✋👌☘🧹☹️🤼‍♂️🚣🏼‍♀️🐌🏚🌉🚝🕔👨‍👩‍👦‍👦🏄🏿‍♂️🙁🚜☎🧟‍♂️🥒🎄🧢📓🧘🏻‍♂️👨‍👩‍👧‍👧🤍🍳🚢🚁⭐🎏🎙🤹🏻‍♀️🤽🏾‍♀️🦭🕸🍫⛪🥻🔎🪚⚖️💷😴😲🤶🌰✈☔🖨🛗👩‍✈️😧🐜✨👩‍🎓🧗🏿‍♂️🐂🦘🚕🎉🎨🏌🏻‍♂️🚣🏻‍♂️🙇🪳🎬🪠🌾⌛📸👨‍🏭🙆‍♀️🤷‍♂️🍹📅🚶‍♀️🦻🦒🍬🏊🏻‍♂️🎊🤢😓🤝🧟🚧🚀🌀🤸🏻‍♀️🏄🏻‍♀️🏄🏾‍♂️")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}