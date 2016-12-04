package main

import (
	"fmt"
	"os"
	"os/exec"
)

func SttyEcho() (err error) {

	cmd := exec.Command("stty", "echo")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("stty echo: %v", err)
	}
	return
}

func SttyNoEcho() (err error) {

	cmd := exec.Command("stty", "-echo")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("stty -echo: %v", err)
	}
	return
}

func SttyIcanon() (err error) {

	cmd := exec.Command("stty", "icanon")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("stty icanon: %v", err)
	}
	return
}

func SttyNoIcanon() (err error) {

	cmd := exec.Command("stty", "-icanon")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("stty -icanon: %v", err)
	}
	return
}
