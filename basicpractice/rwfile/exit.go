package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func exit() {
	dat, err := ioutil.ReadFile("main.go")
	check(err)
	fmt.Print(string(dat))

	werr := ioutil.WriteFile("/tmp/dat1", dat, 0644)
	check(werr)

	f, err := os.Open("main.go")
	check(err)
	defer f.Close()

	wf, err := os.Create("/tmp/dat2")
	check(err)
	defer wf.Close()

	b1 := make([]byte, 5)
	n1, err := f.Read(b1)
	check(err)
	fmt.Printf("%d bytes: %s \n", n1, string(b1))
	wn, err := wf.Write(b1)
	check(err)
	fmt.Printf("%d bytes are write to file", wn)
	wf.Sync()

	o2, err := f.Seek(6, 0)
	check(err)
	b1 = make([]byte, 5)
	n1, err = f.Read(b1)
	check(err)
	fmt.Printf("%d bytes read @ %d: %s \n", n1, o2, string(b1))

	o3, err := f.Seek(6, 0)
	check(err)
	b1 = make([]byte, 5)
	n1, err = f.Read(b1)
	check(err)
	fmt.Printf("%d bytes read @ %d: %s \n", n1, o3, string(b1))

	r4 := bufio.NewReader(f)
	b4, err := r4.Peek(5)
	check(err)
	fmt.Printf("5 bytes: %s \n", string(b4))

	w := bufio.NewWriter(wf)
	n4, err := w.WriteString("bufferred \n")
	check(err)
	fmt.Printf("Wrote %d bytes \n", n4)
	w.Flush()
}
