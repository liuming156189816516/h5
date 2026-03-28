package sync_set

import (
	"bytes"
	"fmt"
	"qqgame/baselib"
	"qqgame/baselib/container/set"
)

// 以并发安全map为底层实现的集合类
// 支持并发安全

type syncHashSet struct {
	m baselib.Map
}

// 新建一个并发安全集合
func NewSyncHashSet() *syncHashSet {
	return &syncHashSet{}
}

/*
	@desc: 添加元素至集合
	@params:
		e: 待添加的元素
	@returns:
		bool: 元素已存在，返回flase；元素不存在，返回true；
*/
func (one *syncHashSet) Add(e interface{}) bool {
	//判断元素是否在set中已存在，若不存在，则添加并将键对应的值置为true，若已存在，则返回false
	_,Loaded:=one.m.LoadOrStore(e,true)
	return !Loaded
}

/*
	@desc: 从集合中删除元素
	@params:
		e: 待删除的元素
	@returns:
*/
func (one *syncHashSet) Remove(e interface{}) {
	one.m.Delete(e)
}


/*
	@desc: 清除集合中的所有元素
	@params:
	@returns:
*/
func (one *syncHashSet) Clear() {
	var newM baselib.Map
	one.m=newM
}

/*
	@desc: 判断集合中是否存在某个元素
	@params:
		e: 待判断的元素
	@returns:
		bool: 如果集合中存在元素e，则返回true；若不存在，则返回false
*/
func (one *syncHashSet) IfContains(e interface{}) bool {
	if _,ok:=one.m.Load(e);!ok {
		return false
	}else{
		return true
	}
}

/*
	@desc: 返回集合的元素个数
	@params:
		e: 待判断的元素
	@returns:
		bool: 如果集合中存在元素e，则返回true；若不存在，则返回false
*/
func (one *syncHashSet) Len() int {
	var sum =0
	one.m.Range(func(k,v interface{})bool{
		sum++
		return true
	})
	return sum
}

/*
	@desc: 返回集合是否为空
	@params:
	@returns:
		bool: 如果集合为空（含元素个数为0），返回true;不为空，则返回false
*/
func (one *syncHashSet) Empty() bool {
	if one.Len()==0{
		return true
	}else{
		return false
	}
}

/*
	@desc: 判断集合one与集合other是否完全相同（即包含的元素完全相同）
	@params:
		other： 指向另一个并发安全集合的指针
	@returns:
		bool: 相同，返回true;不相同，则返回false
*/
func (one *syncHashSet) Same(other *syncHashSet) bool {
	if other == nil {
		return false
	}
	var state bool=true
	oneImage:=set.NewHashSet()
	otherImage:=set.NewHashSet()
	one.m.Range(func(setk,setv interface{})bool{
		if state==true {
			other.m.Range(func(otherkey, othervalue interface{}) bool {
				otherImage.Add(otherkey)
				return true
			})
			state = false
		}
		oneImage.Add(setk)
		return true
	})
	return oneImage.Same(otherImage)
}


/*
	@desc: 返回集合中的所有元素，并以切片的形式返回
	@params:
	@returns:
		[]interface{}: 集合元素组成的切片
*/
func (set *syncHashSet) Elements() []interface{} {
	initialLen := set.Len()
	snapshot := make([]interface{}, initialLen)
	actualLen := 0
	set.m.Range(func(k,v interface{})bool{
		if actualLen < initialLen {
			snapshot[actualLen] = k
		} else {
			snapshot = append(snapshot, k)
		}
		actualLen++
		return true
	})
	if actualLen < initialLen {
		snapshot = snapshot[:actualLen]
	}
	return snapshot
}

/*
	@desc: 将集合转换为字符串，字符串形式为HashSet{A B C}
	@params:
	@returns:
		string: 返回的字符串
*/
func (set *syncHashSet) String() string {
	var buf bytes.Buffer
	buf.WriteString("HashSet{")
	first := true
	set.m.Range(func(k,v interface{})bool{
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%v", k))
		return true
	})
	buf.WriteString("}")
	return buf.String()
}
