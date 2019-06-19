package main

import (
	"fmt"
	"net/url"
)

func exit() {
	s := "http://user:pass@baidu.com:8080/fre/seg?a=1&b=2#fed"
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	fmt.Println(u.Scheme)
	fmt.Println(u.Path)     // /seg
	fmt.Println(u.Fragment) // str after #

	fmt.Println(u.RawQuery) //a=1&b=2
	m, _ := url.ParseQuery(u.RawPath)
	fmt.Println(m, m["a"], m["b"])
}
