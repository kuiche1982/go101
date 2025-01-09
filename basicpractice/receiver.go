package main

type A struct {
  I int
}

func(a A) Set(a int) {
  // a.A is not changed, it's a copy or the original object A
  // Do not expecte changes here affect the original object A
  a.I = a
}

func(a *A) SetWithPointerReceiver(a int) {
  a.I = a
}


func main() {
  obj := A{}
  obj.Set(1)
  fmt.Println(obj.I) // output 0, since Set method is working on a copy of obj
  obj.SetWithPointerReceiver(1)
  fmt.Println(obj.I) // output 1
}
