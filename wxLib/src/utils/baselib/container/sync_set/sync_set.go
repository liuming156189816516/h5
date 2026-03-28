package sync_set

import "qqgame/baselib/container/set"

// 实现两个集合间的并、交、差等操作
// 支持并发安全

//取得内容镜像
func getImage(one *syncHashSet, other *syncHashSet)(*set.HashSet,*set.HashSet){
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
	return oneImage,otherImage
}


// 判断集合 one 是否是集合 other 的超集（other中的元素one中都包含）
func IsSuperset(one *syncHashSet, other *syncHashSet) bool {
	if one == nil || other == nil {
		return false
	}
	oneImage,otherImage:=getImage(one,other)
	oneLen := oneImage.Len()
	otherLen := otherImage.Len()
	if oneLen == 0 || oneLen == otherLen {
		return false
	}
	if oneLen > 0 && otherLen == 0 {
		return true
	}
	for _, v := range other.Elements() {
		if !one.IfContains(v) {
			return false
		}
	}
	return true
}

// 生成集合 one 和集合 other 的并集（one和other中所有元素组成的集合）
func Union(one *syncHashSet, other *syncHashSet) *syncHashSet {
	if one == nil || other == nil {
		return nil
	}
	oneImage,otherImage:=getImage(one,other)
	unionedSet := NewSyncHashSet()
	for _, v := range oneImage.Elements() {
		unionedSet.Add(v)
	}
	if other.Len() == 0 {
		return unionedSet
	}
	for _, v := range otherImage.Elements() {
		unionedSet.Add(v)
	}
	return unionedSet
}

// 生成集合 one 和集合 other 的交集（one和other中共有元素组成的集合）
func Intersect(one *syncHashSet, other *syncHashSet) *syncHashSet {
	if one == nil || other == nil {
		return nil
	}
	oneImage,otherImage:=getImage(one,other)
	intersectedSet := NewSyncHashSet()
	if otherImage.Len() == 0 {
		return intersectedSet
	}
	if oneImage.Len() < otherImage.Len() {
		for _, v := range oneImage.Elements() {
			if otherImage.IfContains(v) {
				intersectedSet.Add(v)
			}
		}
	} else {
		for _, v := range otherImage.Elements() {
			if oneImage.IfContains(v) {
				intersectedSet.Add(v)
			}
		}
	}
	return intersectedSet
}

// 生成集合 one 对集合 other 的差集（one有但other没有的元素组成的集合）
func Difference(one *syncHashSet, other *syncHashSet) *syncHashSet {
	if one == nil || other == nil {
		return nil
	}
	oneImage,otherImage:=getImage(one,other)
	differencedSet := NewSyncHashSet()
	if otherImage.Len() == 0 {
		for _, v := range oneImage.Elements() {
			differencedSet.Add(v)
		}
		return differencedSet
	}
	for _, v := range oneImage.Elements() {
		if !otherImage.IfContains(v) {
			differencedSet.Add(v)
		}
	}
	return differencedSet
}

// 生成集合 one 和集合 other 的对称差集（one和other不共有的元素的集合）
func SymmetricDifference(one *syncHashSet, other *syncHashSet) *syncHashSet {
	if one == nil || other == nil {
		return nil
	}
	oneImage,otherImage:=getImage(one,other)
	SymDifferencedSet := NewSyncHashSet()
	diffA := set.Difference(oneImage, otherImage)
	if other.Len() == 0 {
		for _,val:=range diffA.Elements(){
			SymDifferencedSet.Add(val)
		}
		return SymDifferencedSet
	}
	diffB := set.Difference(otherImage,oneImage)
	for _,val:=range set.Union(diffA, diffB).Elements(){
		SymDifferencedSet.Add(val)
	}
	return SymDifferencedSet
}

