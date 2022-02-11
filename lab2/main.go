package main

import (
	"fmt"
	"goLang/lab2/lru"
	"strconv"
)

func main() {
	testlru := lru.NewCache(3)

	for i := 0; i < 4; i++ {
		err := testlru.Put("key"+strconv.Itoa(i), "val"+strconv.Itoa(i))
		if err != nil {
			fmt.Println(err)
		}
	}
	// changing the value of key 0 to val5
	err := testlru.Put("key"+strconv.Itoa(0), "val"+strconv.Itoa(5))
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
	}
	// get key3 now
	if val, err := testlru.Get("key3"); err != nil {
		fmt.Println(err)
	} else if val.(string) != "val3" {
		fmt.Println(err)
	}
	// get key1 and check for val3 just to throw an error. will return nil
	if val, err := testlru.Get("key1"); err != nil {
		fmt.Println(err)
	} else if val.(string) != "val3" {
		fmt.Println(err)
	}
}
