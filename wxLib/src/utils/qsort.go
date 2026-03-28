package utils
//
//import "sort"
//type qsort_struct struct {
//	arr *[]interface{}
//	less func(interface{}, interface{})bool
//}
////用于给qsort回调
//func (s *qsort_struct) Len() int {
//	return len(*s.arr)
//}
//
////Less():成绩将有低到高排序
//func (s *qsort_struct) Less(i, j int) bool {
//	if s == nil || s.less == nil|| s.arr == nil || i< 0 || i >=len(*s.arr) || j < 0 || j >= len(*s.arr){
//		return false
//	}
//	return s.less((*s.arr)[i],(*s.arr)[j])
//}
//
////Swap()
//func (s *qsort_struct) Swap(i, j int) {
//	s.arr[i], s.arr[j] = s.arr[j], s.arr[i]
//}
//func QSort(arr []interface{}, less func(interface{}, interface{})bool){
//	qs :=&qsort_struct{arr:&arr, less:less}
//	sort.Sort(qs)
//}
//
//
//type Interface interface {
//	// Len is the number of elements in the collection.
//	Len() int
//
//	// Swap swaps the elements with indexes i and j.
//	Swap(i, j int)
//}
//
//func QRandom( arr Interface){
//	if arr == nil{
//		return
//	}
//	alen:=arr.Len()
//	if alen<=1{
//		return
//	}
//
//	for i:=0; i < alen; i++{
//		j := RandGen(0, alen-1)
//		arr.Swap(i, j)
//	}
//
//}
//
