package utils

import (
	"reflect"
)

// CheckIsNil 这个函数检查一个对象是否是空
//
// 它会先检查对象是否实现了 NilAble 这个接口
// 如果没有实现就走反射 反射如果不支持就 直接 return false
// 可能会有误差
func CheckIsNil(val any) bool {
	if val == nil {
		return true
	}
	type NilAble interface {
		IsNil() bool
	}
	// 先调用自身的方法判断是不是 nil
	if v, ok := val.(NilAble); ok {
		return v.IsNil()
	} else {
		// 自身的方法不行 就走反射
		vo := reflect.ValueOf(val)
		k := vo.Kind()
		switch k {
		case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer,
			reflect.UnsafePointer, reflect.Interface, reflect.Slice:
			return vo.IsNil()
		default:
			return false
		}
	}
}
