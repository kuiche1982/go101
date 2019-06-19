package main

import (
	"fmt"
	"io/ioutil"
	data "kuitest/additional/embedresource/resourcepath"
)

//  //go:generate go run -tags=dev assets_generate.go
//go:generate vfsgendev -source="kuitest/additional/embedresource/resourcepath".Assets
func main() {
	//https://tech.townsourced.com/post/embedding-static-files-in-go/
	//https://github.com/shurcooL/vfsgen
	//go get -u github.com/shurcooL/vfsgen/cmd/vfsgendev
	// this is something not work. but it's good idea to generate comile the resoruce and use the generated assets_vfsdata.go.
	fs := data.Assets
	f, err := fs.Open("/a.pptx")
	if err != nil {
		s, err := ioutil.ReadAll(f)
		fmt.Println(string(s), err)
	}
}
