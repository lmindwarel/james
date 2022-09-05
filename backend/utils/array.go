package utils

import "reflect"

func InArray(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if reflect.DeepEqual(arr.Index(i).Interface(), item) {
			return true
		}
	}

	return false
}
