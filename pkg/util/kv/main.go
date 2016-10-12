package main

import (
	"fmt"
	cache "github.com/dadgar/onecache/ttlstore"
	"log"
	"os"
	"time"
)

var (
	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func main() {
	cache, err := cache.New(10000, logger)

	if err != nil {
		log.Println(err)
	}

	t := time.Now().Add(1 * time.Second).Unix()
	//
	cache.Add("123", []byte("456"), t, 99)

	fmt.Println(cache.List())

	time.Sleep(time.Second)

	fmt.Println(len(cache.List()))
	//
	//for {
	//	v, err := cache.Get("123")
	//	if err != nil {
	//		log.Println("----", err)
	//	}
	//
	//	log.Println(v)
	//	time.Sleep(time.Second * 1)
	//
	//	cache.Set("123", []byte("456"), t, 99)
	//}

}
