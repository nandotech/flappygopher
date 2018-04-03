package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type pipes struct {
	mu      sync.RWMutex
	texture *sdl.Texture
	speed   int32

	pipes []*pipe
}

func (ps *pipes) touch(b *bird) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	for _, p := range ps.pipes {
		p.touch(b)
	}
}

func newPipes(r *sdl.Renderer) (*pipes, error) {
	p, err := img.LoadTexture(r, "res/imgs/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("could not load pipe image: %v", err)
	}
	ps := &pipes{texture: p, speed: 2}
	rand.Seed(time.Now().UnixNano())
	go func() {
		for {
			ps.mu.Lock()
			ps.pipes = append(ps.pipes, newPipe())
			ps.mu.Unlock()
			time.Sleep(time.Duration(1) * time.Second)
		}
	}()
	return ps, nil
}

func (ps *pipes) paint(r *sdl.Renderer) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		if err := p.paint(r, ps.texture); err != nil {
			return err
		}
	}
	return nil
}
func (ps *pipes) update() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var rem []*pipe

	for _, p := range ps.pipes {
		p.mu.Lock()
		p.x -= ps.speed
		p.mu.Unlock()
		if p.x+p.w > 0 {
			rem = append(rem, p)
		}
	}
	ps.pipes = rem
}

func (ps *pipes) restart() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.pipes = nil
}

func (ps *pipes) destroy() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.texture.Destroy()
}

type pipe struct {
	mu sync.RWMutex

	x int32
	h int32
	w int32

	inverted bool
}

func newPipe() *pipe {
	rand.Seed(time.Now().UnixNano())
	return &pipe{
		x:        800,
		h:        rand.Int31n(250) + 100,
		w:        50,
		inverted: rand.Float32() < 0.5,
	}
}

func (p *pipe) paint(r *sdl.Renderer, texture *sdl.Texture) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	flip := sdl.FLIP_NONE
	rect := &sdl.Rect{X: p.x, Y: 600 - p.h, W: p.w, H: p.h}
	if p.inverted {
		rect.Y = 0
		flip = sdl.FLIP_VERTICAL
	}

	if err := r.CopyEx(texture, nil, rect, 0, nil, flip); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	return nil
}

func (p *pipe) restart() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.x = 400
}

func (p *pipe) touch(b *bird) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	b.touch(p)
}
