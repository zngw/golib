// @Title
// @Description $
// @Author  55
// @Date  2022/5/30
package main

import (
	"fmt"
	"github.com/zngw/golib/zchan"
	"time"
)

func main() {
	zc, err := zchan.New(4)
	if err != nil {
		panic(err.Error())
	}

	go func() {
		// 写入channel数据
		// 10毫秒写入1次
		for i := 1000; i < 2000; i++ {
			zc.In <- i
			fmt.Printf("写入数据：%v, chan长度：%v， Buf长度： %v \n", i, zc.Len(), zc.BufLen())
			time.Sleep(time.Millisecond)
		}
	}()

	go func() {
		// 写入channel数据
		// 10毫秒写入1次
		for i := 0; i < 1000; i++ {
			zc.In <- i
			fmt.Printf("写入数据：%v, chan长度：%v， Buf长度： %v \n", i, zc.Len(), zc.BufLen())
			time.Sleep(time.Millisecond)
		}
	}()

	for v := range zc.Out {
		// 20 毫毛读取一次数据
		fmt.Printf("读取入数据：%v, chan长度：%v， Buf长度：%v \n", v, zc.Len(), zc.BufLen())
		time.Sleep(20 * time.Millisecond)
	}
}
