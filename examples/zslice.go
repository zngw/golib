package main

import (
	"fmt"

	"github.com/zngw/golib/zslice"
)

func main() {
	// 示例：存储字符串
	slice := zslice.NewSlice[string]()

	// 添加数据
	for _, v := range []string{"aa", "bb", "cc", "aa", "ddd"} {
		slice.Append(v)
	}

	fmt.Println("初始数据:", slice.ToSlice())

	// 排序
	slice.Sort(func(a, b string) bool {
		return a < b
	})
	fmt.Println("排序后:", slice.ToSlice())

	// 查找
	if val, ok := slice.Find(func(s string) bool { return s == "bb" }); ok {
		fmt.Println("找到:", val)
	}

	// 查找所有 "aa"
	aa := slice.FindAll(func(s string) bool { return s == "aa" })
	fmt.Println("所有 aa:", aa)

	// 按值删除
	count := slice.DeleteBy(func(a string) bool { return a == "aa" })
	fmt.Printf("删除了 %d 个 'aa'\n", count)
	fmt.Println("删除后:", slice.ToSlice())

	// 统计长度
	fmt.Println("当前长度:", slice.Len())

	// 统计满足条件的元素
	longCount := slice.Count(func(s string) bool { return len(s) > 2 })
	fmt.Printf("长度大于2的元素个数: %d\n", longCount)
}
