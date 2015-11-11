package main

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

func execCommand(cmd string, wg *sync.WaitGroup) {
	parts := strings.Fields(cmd)
	command := parts[0]
	args := parts[1:len(parts)]

	out, err := exec.Command(command, args...).CombinedOutput()

	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", string(out))
	wg.Done()
}
