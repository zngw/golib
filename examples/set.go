package main

import (
	"fmt"
	"github.com/zngw/golib/set"
)

func main() {
	s := set.New()

	s.Add(1)
	s.Add(1)
	s.Add(0)
	s.Add(2)
	s.Add(4)
	s.Add(3)
	fmt.Println("输出数据：", s.List())

	s.Clear()
	if s.IsEmpty() {
		fmt.Println("空Set")
	}

	s.Add(1)
	s.Add(2)
	s.Add(3)

	if s.Has(2) {
		fmt.Println("2 存在")
	}

	s.Remove(2)
	s.Remove(3)
	fmt.Println("无序的切片", s.List())
}
