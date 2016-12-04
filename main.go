package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {

	var st State

	flag.Parse()
	for i, arg := range flag.Args() {
		var m, s int
		if _, err := fmt.Sscanf(arg, "%d/%d", &m, &s); err != nil {
			log.Fatalf("%v: %v", arg, err)
		}
		l := &Laser{Id: i + 1, Meters: m, Seconds: s}
		st.Lasers = append(st.Lasers, l)
	}

	if err := SttyNoEcho(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := SttyEcho(); err != nil {
			log.Println(err)
		}
	}()

	if err := SttyNoIcanon(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := SttyIcanon(); err != nil {
			log.Println(err)
		}
	}()

	input := make(chan byte, 32)
	go func() {
		buf := make([]byte, 4)
		for {
			n, err := os.Stdin.Read(buf)
			for _, b := range buf[0:n] {
				input <- b
			}
			if err != nil {
				if err != io.EOF {
					log.Println(err)
				}
				break
			}
		}
		close(input)
	}()

	var text string
	var tlen int
	{
		t := st.Text()
		text = t.Text()
		tlen = t.StringLen()
	}
	fmt.Print(text)
	tick := time.Tick(100 * time.Millisecond)

	for input != nil {

		select {
		case <-tick:
		case b, ok := <-input:
			if !ok || b == 4 { // EOT
				input = nil
				break
			}
			if b == '+' && st.Plus() {
			} else if b == '-' && st.Minus() {
			} else if (b == 'k' || b == 'h') && st.Left() {
			} else if (b == 'j' || b == 'l') && st.Right() {
			} else if b == 127 && st.Backspace() {
			} else if b >= '0' && b <= '9' && st.Digit(b, time.Now()) {
			} else if o := Ores[b]; o != nil && st.Ore(o) {
			}
		}

		for i := range st.Tick(time.Now()) {
			title := fmt.Sprintf("laser %v", i.Id)
			body := fmt.Sprintf("%v depleted", i.Lock.Ore.Name)
			exec.Command("notify-send", title, body).Run()
		}

		t := st.Text()
		s := t.Text()
		if s != text {
			ss := "\r" + s
			l := t.StringLen()
			if l < tlen {
				ss += strings.Repeat(" ", tlen-l)
				ss += "\r" + s
			}
			fmt.Print(ss)
			text, tlen = s, l
		}
	}
	fmt.Println()
}
