package main

import (
	"fmt"
	"time"
)

type Ore struct {
	Unit   string
	Name   string
	Volume int
}

// http://wiki.eveuniversity.org/Asteroids_and_Ore#Ore_Chart
var Ores = map[byte]*Ore{
	'v': &Ore{"v", "Veldspar", 10},
	's': &Ore{"s", "Scordite", 15},
	'y': &Ore{"py", "Pyroxeres", 30},
	'l': &Ore{"pl", "Plagioclase", 35},
	'o': &Ore{"o", "Omber", 60},
	'k': &Ore{"k", "Kernite", 120},
}

type Laser struct {
	Id      int
	Meters  int
	Seconds int
	Lock    *Asteroid
	Last    time.Time
}

type Asteroid struct {
	Units int
	Ore   *Ore
	Locks map[*Laser]struct{}
}

func (a *Asteroid) Text() (t Texts) {
	if a.Units <= 0 {
		t = append(t, Red)
	} else if len(a.Locks) > 0 {
		t = append(t, Green)
	}
	t = append(t, String(fmt.Sprintf("%v%v", a.Units, a.Ore.Unit)))
	t = append(t, Reset)
	if len(a.Locks) > 0 {
		t = append(t, String("["))
		first := true
		for l := range a.Locks {
			if !first {
				t = append(t, String(","))
			}
			t = append(t, String(fmt.Sprintf("%v", l.Id)))
			first = false
		}
		t = append(t, String("]"))
	}
	return
}

type State struct {
	Lasers      []*Laser
	Asteroids   []*Asteroid
	SelAsteroid int
	NewAsteroid *string
}

func (st *State) Text() (t Texts) {
	for i, a := range st.Asteroids {
		if i != 0 {
			t = append(t, String(", "))
		}
		if i == st.SelAsteroid {
			t = append(t, Blue)
			t = append(t, String("*"))
			t = append(t, Reset)
		}
		t = append(t, a.Text())
	}
	t = append(t, String("> "))
	if st.NewAsteroid != nil {
		t = append(t, String("+"))
		t = append(t, String(*st.NewAsteroid))
	}
	return
}
