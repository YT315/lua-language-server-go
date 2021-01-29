package auxiliary

import (
	"fmt"
	"lualsp/syntax"
	"reflect"
	"strings"
)

//Traversal 遍历抽象语法树，绘制出树结构
func Traversal(node interface{}, vist func(syntax.Node)) {
	v := reflect.ValueOf(node)
	if v.Kind() == reflect.Invalid {
		fmt.Printf("nill\n")
		return
	}
	if n, ok := node.(syntax.Node); ok {
		vist(n)
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t := v.Type()
		//fmt.Printf("\n")
		for i := 0; i < t.NumField(); i++ {
			//基类不进
			nm := t.Field(i).Name
			if strings.Index(nm, "Base") != -1 {
				continue
			}
			sv := v.Field(i)
			switch sv.Kind() {
			case reflect.Slice:
				for j := 0; j < sv.Len(); j++ {
					ifc := sv.Index(j).Interface()
					if _, ok := ifc.(syntax.Node); ok {
						Traversal(ifc, vist)
					}
				}
			case reflect.Ptr:
				ifc := sv.Interface()
				if _, ok := ifc.(syntax.Node); ok {
					Traversal(ifc, vist)
				}
			case reflect.Interface:
				ifc := sv.Interface()
				if _, ok := ifc.(syntax.Node); ok {
					Traversal(ifc, vist)
				}
			default:
			}
		}
	} else {
		panic(v.Kind())
	}

}
