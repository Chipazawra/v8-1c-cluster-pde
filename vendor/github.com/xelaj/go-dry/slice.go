package dry

import (
	"fmt"
	"reflect"
)

func SliceIndex(i, item interface{}) int {
	ival := reflect.ValueOf(i)
	if ival.Type().Kind() != reflect.Slice {
		panic("not a slice: " + ival.Type().String())
	}
	if ival.Type().Elem().String() != reflect.TypeOf(item).String() {
		panic("different types of slice and item")
	}

	for i := 0; i < ival.Len(); i++ {
		if reflect.DeepEqual(ival.Index(i).Interface(), item) {
			return i
		}
	}

	return -1
}

func DeleteIndex(slice interface{}, i int) interface{} {
	return cutSliceWithMode(slice, i, i+1, true)
}

func SliceCut(slice interface{}, i, j int) interface{} {
	return cutSliceWithMode(slice, i, j, false)
}

func cutSliceWithMode(slice interface{}, i, j int, deleteInsteadCut bool) interface{} {
	panicIndexStr := fmt.Sprintf("[%v:%v]", i, j)
	if deleteInsteadCut {
		panicIndexStr = fmt.Sprintf("[%v]", i)
	}

	ival := reflect.ValueOf(slice)
	if ival.Type().Kind() != reflect.Slice {
		panic("not a slice: " + ival.Type().String())
	}

	if i > j {
		panic("end less than start " + panicIndexStr)
	}
	if i < 0 || j < 0 {
		panic("slice index " + panicIndexStr + " out of bounds")
	}
	if ival.Len()-1 < i || ival.Len() < j {
		panic(fmt.Sprintf("index out of range %v with length %v", panicIndexStr, ival.Len()))
	}

	return reflect.AppendSlice(ival.Slice(0, i), ival.Slice(j, ival.Len())).Interface()
}

func SliceExpand(slice interface{}, i, j int) interface{} {
	panicIndexStr := fmt.Sprintf("[%v]", i)

	ival := reflect.ValueOf(slice)
	if ival.Type().Kind() != reflect.Slice {
		panic("not a slice: " + ival.Type().String())
	}

	if i < 0 {
		panic("slice index " + panicIndexStr + " out of bounds")
	}
	if j < 0 {
		panic(fmt.Sprintf("can't expand slice on %v points", j))
	}
	if ival.Len()-1 < i {
		panic(fmt.Sprintf("index out of range %v with length %v", panicIndexStr, ival.Len()))
	}

	zeroitems := reflect.MakeSlice(ival.Type(), j, j)
	part := reflect.AppendSlice(zeroitems, ival.Slice(i, ival.Len()))
	return reflect.AppendSlice(ival.Slice(0, i), part).Interface()
}

func SliceToInterfaceSlice(in interface{}) []interface{} {
	ival := reflect.ValueOf(in)
	if ival.Type().Kind() != reflect.Slice {
		panic("not a slice: " + ival.Type().String())
	}

	res := make([]interface{}, ival.Len())

	for i := 0; i < ival.Len(); i++ {
		res[i] = ival.Index(i).Interface()
	}
	return res
}

func MapKeys(in interface{}) interface{} {
	ival := reflect.ValueOf(in)
	if ival.Type().Kind() != reflect.Map {
		panic("not a map: " + ival.Type().String())
	}

	keys := ival.MapKeys()

	items := reflect.MakeSlice(reflect.SliceOf(ival.Type().Key()), len(keys), len(keys))
	for i, key := range keys {
		items.Index(i).Set(key)
	}

	return items.Interface()
}

func SliceUnique(in interface{}) interface{} {
	ival := reflect.ValueOf(in)
	if ival.Type().Kind() != reflect.Slice {
		panic("not a slice: " + ival.Type().String())
	}

	res := reflect.MakeMap(reflect.MapOf(ival.Type().Elem(), reflect.TypeOf(struct{}{})))

	for i := 0; i < ival.Len(); i++ {
		res.SetMapIndex(ival.Index(i), reflect.ValueOf(struct{}{}))
	}
	return res.Interface()
}
