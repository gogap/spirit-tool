package main

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func execCommand(cmd string) (out []byte, err error) {
	parts := strings.Fields(cmd)
	command := parts[0]
	args := parts[1:len(parts)]

	out, err = exec.Command(command, args...).CombinedOutput()

	return
}

func execCommandWithDir(cmd string, dir string) (out []byte, err error) {
	parts := strings.Fields(cmd)
	command := parts[0]
	args := parts[1:len(parts)]

	cmder := exec.Command(command, args...)
	cmder.Dir = dir

	out, err = cmder.CombinedOutput()

	return
}

func execute(cmd string, dir string, bindSTD bool, envs []string) (cmder *exec.Cmd, err error) {
	parts := strings.Fields(cmd)
	command := parts[0]
	args := parts[1:len(parts)]

	commander := exec.Command(command, args...)

	if dir != "" {
		commander.Dir = dir
	}

	if bindSTD {
		commander.Stderr = os.Stderr
		commander.Stdout = os.Stdout
		commander.Stdin = os.Stdin
	}

	commander.Env = append(os.Environ(), envs...)

	if err = commander.Start(); err != nil {
		return
	}

	cmder = commander

	return
}

func killProcess(pid int) (err error) {
	err = syscall.Kill(pid, syscall.SIGKILL)
	return
}

func stopProcess(pid int) (err error) {
	err = syscall.Kill(pid, syscall.SIGTERM)
	return
}

func isProcessAlive() {

}
