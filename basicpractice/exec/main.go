package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func mainx() {
	binary, lookError := exec.LookPath("ls")
	if lookError != nil {
		panic(lookError)
	}

	args := []string{"ls", "-a", "-l", "-h"}
	env := os.Environ()

	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}

func main() {
	cmd := exec.Command("cat", "main.go")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	_, err = stdin.Write([]byte(""))
	if err != nil {
		panic(err)
	}
	stdin.Close()

	//cmd.StdoutPipe()

	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	cmd.Start()
}
