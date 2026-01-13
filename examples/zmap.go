package main

import (
	"fmt"

	"github.com/zngw/golib/zmap"
)

type User struct {
	Name  string
	Score int
}

func main() {
	userCache := zmap.New[string, *User]()

	// Store
	userCache.Store("u1", &User{Name: "Zngw", Score: 100})

	// Load
	u, ok := userCache.Load("u1")
	if ok {
		fmt.Println("User:", u.Name)
	}

	// LoadOrStore
	userCache.LoadOrStore("u2", &User{Name: "Guoke", Score: 90})

	// Range
	userCache.Range(func(id string, u *User) bool {
		fmt.Printf("ID: %s, Name: %s, Score: %d\n", id, u.Name, u.Score)
		return true // continue
	})

	// CAS
	ok = userCache.CompareAndSwap("u1", u, &User{
		Name:  "Zngw",
		Score: 80,
	})
	fmt.Println(ok)

	// 再次遍历输出Range
	userCache.Range(func(id string, u *User) bool {
		fmt.Printf("ID: %s, Name: %s, Score: %d\n", id, u.Name, u.Score)
		return true // continue
	})

	// Get all keys
	keys := userCache.Keys()
	fmt.Println("All user IDs:", keys)
}
