package core

import "time"

type FPS struct {
	lastUpdate time.Time
	cnt        int
	fps        int
}

func (f *FPS) Update() {
	f.cnt++
	now := time.Now()
	p := now.Sub(f.lastUpdate)
	if p >= time.Second {
		f.fps = int(float64(f.cnt) / p.Seconds())
		f.cnt = 0
		f.lastUpdate = now
	}
}

func (f *FPS) Fps() int {
	return f.fps
}
