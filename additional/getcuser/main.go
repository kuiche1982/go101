package main

import (
	"fmt"
	"os/user"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic("unable to get the current user required for UCM SDK, err: " + err.Error())
	}
	fmt.Println(usr)
}
