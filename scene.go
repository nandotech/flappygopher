package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	time int

	bg    *sdl.Texture
	birds *bird
	pipes *pipes
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}
	b, err := newBird(r)
	if err != nil {
		return nil, err
	}

	p, err := newPipes(r)
	if err != nil {
		return nil, err
	}
	return &scene{bg: bg, birds: b, pipes: p}, nil
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)

		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			//	log.Printf("event: %T", e)
			case <-tick:
				s.update()
				if s.birds.isDead() {
					drawTitle(r, "Game Over")
					time.Sleep(1 * time.Second)
					s.restart()
				}
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}

	}()
	return errc
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.KeyboardEvent:
		s.birds.jump()
		return false
	default:
		log.Printf("unknown event: %T", event)
		return false
	}
}

func (s *scene) paint(r *sdl.Renderer) error {
	s.time++
	r.Clear()

	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	if err := s.birds.paint(r); err != nil {
		return err
	}
	if err := s.pipes.paint(r); err != nil {
		return err
	}

	r.Present()
	return nil
}

func (s *scene) update() {
	s.birds.update()
	s.pipes.update()
	s.pipes.touch(s.birds)
}

func (s *scene) restart() {
	s.birds.restart()
	s.pipes.restart()
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.birds.destroy()
	s.pipes.destroy()
}
