package comtools

import (
	"checkManager/database"
	"reflect"
)

func StructArr2TwoArr(obj interface{}) [][]string {
	var datas = make([][]string, 0)
	slice, ok := database.CreateAnyTypeSlice(obj)
	if !ok {
		return datas
	}
	for _, value := range slice {
		v := reflect.ValueOf(value)
		count := v.NumField()
		var data []string
		for i := 1; i < count; i++ { // 不要第一个字段
			data = append(data, v.Field(i).String())
		}
		datas = append(datas, data)
	}
	return datas
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
