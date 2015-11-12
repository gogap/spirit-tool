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

func execute(cmd string, dir string) (cmder *exec.Cmd, err error) {
	parts := strings.Fields(cmd)
	command := parts[0]
	args := parts[1:len(parts)]

	commander := exec.Command(command, args...)

	commander.Dir = dir
	commander.Stderr = os.Stderr
	commander.Stdin = os.Stdin
	commander.Stdout = os.Stdout

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

func isProcessAlive(pid int) bool {
	p, _ := os.FindProcess(pid)
	if e := p.Signal(syscall.Signal(0)); e == nil {
		return true
	}
	return false
}
