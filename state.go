package main

import (
	"fmt"
	"time"
)

func (st *State) Plus() bool {
	if st.NewAsteroid == nil {
		s := ""
		st.NewAsteroid = &s
		return true
	}
	return false
}

func (st *State) Minus() bool {
	if st.NewAsteroid == nil && len(st.Asteroids) > 0 {
		a := st.Asteroids[st.SelAsteroid]
		for l := range a.Locks {
			st.stop(l, a)
		}
		if len(st.Asteroids) > 0 && a == st.Asteroids[st.SelAsteroid] {
			st.deplete1(st.SelAsteroid)
		}
		return true
	}
	return false
}

func (st *State) Left() bool {
	if st.NewAsteroid == nil && st.SelAsteroid > 0 {
		st.SelAsteroid--
		return true
	}
	return false
}

func (st *State) Right() bool {
	if st.NewAsteroid == nil && st.SelAsteroid < len(st.Asteroids)-1 {
		st.SelAsteroid++
		return true
	}
	return false
}

func (st *State) Digit(b byte, now time.Time) bool {
	if st.NewAsteroid != nil {
		s := *st.NewAsteroid + string(b)
		st.NewAsteroid = &s
		return true
	}
	if len(st.Asteroids) > 0 {
		i, s := int(b-'1'), st.SelAsteroid
		if i < len(st.Lasers) {
			l, a := st.Lasers[i], st.Asteroids[s]
			if l.Lock == nil {
				st.start(l, a, now)
			} else if l.Lock == a {
				st.stop(l, a)
			} else {
				st.stop(l, l.Lock)
				st.start(l, a, now)
			}
		}
	}
	return false
}

func (st *State) Backspace() bool {
	if st.NewAsteroid != nil {
		if len(*st.NewAsteroid) == 0 {
			st.NewAsteroid = nil
		} else {
			s := (*st.NewAsteroid)[0 : len(*st.NewAsteroid)-1]
			st.NewAsteroid = &s
		}
		return true
	}
	return false
}

func (st *State) Ore(o *Ore) bool {
	if st.NewAsteroid != nil && len(*st.NewAsteroid) > 0 {
		var q int
		if _, err := fmt.Sscanf(*st.NewAsteroid, "%d", &q); err == nil {
			a := &Asteroid{q, o, make(map[*Laser]struct{})}
			st.Asteroids = append(st.Asteroids, a)
			st.NewAsteroid = nil
			st.SelAsteroid = len(st.Asteroids) - 1
		}
		return true
	}
	return false
}

func (st *State) start(l *Laser, a *Asteroid, now time.Time) {
	l.Lock, l.Last = a, now
	a.Locks[l] = struct{}{}
}

func (st *State) stop(l *Laser, a *Asteroid) {
	l.Lock, l.Last = nil, time.Time{}
	delete(a.Locks, l)
	if a.Units <= 0 && len(a.Locks) == 0 {
		st.deplete2(a)
	}
}

func (st *State) deplete2(a *Asteroid) {
	for i, t := range st.Asteroids {
		if t == a {
			st.deplete1(i)
			break
		}
	}
}

func (st *State) deplete1(i int) {
	l := len(st.Asteroids)
	copy(st.Asteroids[i:l-1], st.Asteroids[i+1:l])
	st.Asteroids[l-1], st.Asteroids = nil, st.Asteroids[0:l-1]
	if st.SelAsteroid > i {
		st.SelAsteroid--
	} else if st.SelAsteroid == i && i > 0 {
		st.SelAsteroid--
	}
}

func (st *State) Tick(now time.Time) (depleted map[*Laser]struct{}) {
	units := make(map[*Asteroid]int, len(st.Asteroids))
	for _, l := range st.Lasers {
		if a := l.Lock; a != nil {
			tick := time.Duration(l.Seconds) * time.Second *
				time.Duration(a.Ore.Volume) / time.Duration(l.Meters*100)
			ticks := now.Sub(l.Last) / tick
			l.Last = l.Last.Add(ticks * tick)
			units[a] += int(ticks)
		}
	}
	depleted = make(map[*Laser]struct{}, len(st.Lasers))
	for a, u := range units {
		old := a.Units
		a.Units -= u
		if old > 0 && a.Units <= 0 {
			for l := range a.Locks {
				depleted[l] = struct{}{}
			}
		}
	}
	return
}
