package auxiliary

import (
	"fmt"
	"reflect"
	"strings"
)

var deep = 0
var flag []bool = []bool{false}

//绘制树干
func putvid() {
	for i := 0; i < deep; i++ {
		if flag[i+1] {
			fmt.Printf("│  ")
		} else {
			fmt.Printf("   ")
		}
	}
	fmt.Printf("\n")
}

//绘制中间树叉
func putmid() {
	for i := 0; i < deep-1; i++ {
		if flag[i+1] {
			fmt.Printf("│  ")
		} else {
			fmt.Printf("   ")
		}
	}
	fmt.Printf("├──")
}

//绘制最后树叉
func putend() {
	for i := 0; i < deep-1; i++ {
		if flag[i+1] {
			fmt.Printf("│  ")
		} else {
			fmt.Printf("   ")
		}
	}
	fmt.Printf("└──")
}

//DrawTree 递归遍历抽象语法树，绘制出树结构
func DrawTree(node interface{}, arg ...string) {
	v := reflect.ValueOf(node)
	t := reflect.TypeOf(node)
	if len(arg) > 0 {
		fmt.Printf(arg[0])
		fmt.Printf(":")
	}
	if v.Kind() == reflect.Invalid {
		fmt.Printf("nill\n")
		return
	}
	deep++
	if len(flag) <= deep {
		flag = append(flag, true)
	} else {
		flag[deep] = true
	}
	switch v.Kind() {
	case reflect.Slice:
		fmt.Printf(t.String()) //[]xxxx
		fmt.Printf("\n")
		for i := 0; i < v.Len(); i++ {
			putvid()
			if i == v.Len()-1 {

				putend()
				flag[deep] = false
			} else {
				putmid()
			}
			DrawTree(v.Index(i).Interface())
		}

	case reflect.Ptr:
		v = v.Elem()
		t := t.Elem()
		fmt.Println(t.String()) //xxxstruct
		//fmt.Printf("\n")
		for i := 0; i < t.NumField(); i++ {
			if nm := t.Field(i).Name; strings.Index(nm, "Base") == -1 {
				putvid()
				if i == t.NumField()-1 {
					putend()
					flag[deep] = false
				} else {
					putmid()
				}
				DrawTree(v.Field(i).Interface(), nm)
			}
		}

	case reflect.String:
		fmt.Printf(v.Kind().String())
		fmt.Printf("\n")
		putvid()
		putend()
		flag[deep] = false
		fmt.Println(v.String())

	case reflect.Bool:
		fmt.Printf(v.Kind().String())
		fmt.Printf("\n")
		putvid()
		putend()
		flag[deep] = false
		fmt.Println(v.Bool())

	case reflect.Float64:
		fmt.Printf(v.Kind().String())
		fmt.Printf("\n")
		putvid()
		putend()
		flag[deep] = false
		fmt.Println(v.Float())

	default:
		panic(v.Kind())
	}
	deep--
}
