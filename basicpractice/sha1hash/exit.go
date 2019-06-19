package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
)

func exit() {
	s := "sha1 this string"
	h := sha1.New()
	h.Write([]byte(s))

	bs := h.Sum(nil)

	fmt.Println(s)
	fmt.Printf("%X\n", bs)

	bs1 := h.Sum([]byte("salt"))

	fmt.Println(s)
	fmt.Printf("salted %X\n", bs1)

	m := md5.New()
	m.Write([]byte(s))
	ms := m.Sum(nil)
	fmt.Println(s)
	fmt.Printf("%X\n", ms)

	sEnc := base64.StdEncoding.EncodeToString([]byte(s))
	fmt.Println("base64 encoding: ", string(sEnc))

	sDec, err := base64.StdEncoding.DecodeString(sEnc)
	fmt.Println("base64 decoded string and error :", string(sDec), err)

}
