package main

import (
	"encoding/json"
	"fmt"
)

type Ppl struct {
	Name string `json:"pName"`
	Age  int    `json:"pAge"`
}

func p() {
	ppls := []Ppl{
		Ppl{
			Name: "kui",
			Age:  3,
		},
		Ppl{
			Name: "4",
			Age:  4,
		}, Ppl{
			Name: "1",
			Age:  1,
		},
		Ppl{
			Name: "2",
			Age:  2,
		},
	}

	fmt.Println(ppls)
	bolB, err := json.Marshal(ppls)
	fmt.Println(string(bolB), err)

	bolB, err = json.Marshal(ppls[0])
	fmt.Println(string(bolB), err)

	um := []byte(`[{"pName":"kui", "pAge":33},{"pName":"kui+1", "pAge":34}]`)
	var o []Ppl
	json.Unmarshal(um, &o)
	fmt.Println(o)
}
