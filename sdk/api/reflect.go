package api

import "reflect"

// CopyProperties 根据目标对象拷贝源对象的值,传递的值要用指针！！
func CopyProperties(from, to interface{}) {
	fType := reflect.TypeOf(from)
	fValue := reflect.ValueOf(from)
	if fType.Kind() == reflect.Ptr {
		fType = fType.Elem()
		fValue = fValue.Elem()
	}
	tType := reflect.TypeOf(to)
	tValue := reflect.ValueOf(to)
	if tType.Kind() == reflect.Ptr {
		tType = tType.Elem()
		tValue = tValue.Elem()
	}
	for i := 0; i < tType.NumField(); i++ {
		for j := 0; j < fType.NumField(); j++ {
			if tType.Field(i).Name == fType.Field(j).Name &&
				fType.Field(j).Type.ConvertibleTo(tType.Field(i).Type) {
				tValue.Field(i).Set(fValue.Field(j))
			}
		}
	}
}
