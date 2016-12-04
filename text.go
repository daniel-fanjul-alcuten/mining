package main

import (
	"log"
	"os"
	"os/exec"
)

type Text interface {
	Text() string
	String() string
	TextLen() int
	StringLen() int
}

type Color string

func (c Color) Text() string {
	return string(c)
}

func (c Color) String() string {
	return ""
}

func (c Color) TextLen() int {
	return len(c)
}

func (c Color) StringLen() int {
	return 0
}

var Black, Red, Green, Yellow, Blue, Magenta, Cyan, White, Reset Color

func init() {
	log.SetFlags(0)
	Black = tput("setaf", "0")
	Red = tput("setaf", "1")
	Green = tput("setaf", "2")
	Yellow = tput("setaf", "3")
	Blue = tput("setaf", "4")
	Magenta = tput("setaf", "5")
	Cyan = tput("setaf", "6")
	White = tput("setaf", "7")
	Reset = tput("sgr0")
}

func tput(args ...string) Color {
	cmd := exec.Command("tput", args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	bytes, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return Color("")
	}
	return Color(bytes)
}

type String string

func (s String) Text() string {
	return string(s)
}

func (s String) String() string {
	return string(s)
}

func (s String) TextLen() int {
	return len(s)
}

func (s String) StringLen() int {
	return len(s)
}

type Texts []Text

func (tt Texts) Text() (s string) {
	for _, t := range tt {
		s += t.Text()
	}
	return
}

func (tt Texts) String() (s string) {
	for _, t := range tt {
		s += t.String()
	}
	return
}

func (tt Texts) TextLen() (l int) {
	for _, t := range tt {
		l += t.TextLen()
	}
	return
}

func (tt Texts) StringLen() (l int) {
	for _, t := range tt {
		l += t.StringLen()
	}
	return
}
