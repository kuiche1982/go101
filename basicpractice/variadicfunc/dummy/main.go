package main

import (
	"fmt"
)

type aa struct {
	aa int
	bb string
}

func main() {
	a := &aa{
		aa: 1,
		bb: "somestring",
	}
	b := *a
	a.aa = 2
	fmt.Println(a)
	fmt.Println(b)
}

// func BuildTypeAtRuntime() {
// 	typ := reflect.StructOf([]reflect.StructField{
// 		{
// 			Name: "Height",
// 			Type: reflect.TypeOf(int(0)),
// 			Tag:  `json:"height"`,
// 		},
// 		{
// 			Name: "Age",
// 			Type: reflect.TypeOf(int(0)),
// 			Tag:  `json:"age"`,
// 		},
// 	})

// 	v := reflect.New(typ).Elem()
// 	v.Field(0).SetFloat(0.4)
// 	v.Field(1).SetInt(2)
// 	s := v.Addr().Interface()

// 	w := new(bytes.Buffer)
// 	if err := json.NewEncoder(w).Encode(s); err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("value: %+v\n", s)
// 	fmt.Printf("json:  %s", w.Bytes())

// 	r := bytes.NewReader([]byte(`{"height":1.5,"age":10}`))
// 	if err := json.NewDecoder(r).Decode(s); err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("value: %+v\n", s)
// }
