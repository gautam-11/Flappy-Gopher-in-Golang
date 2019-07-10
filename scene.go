package main

import (
	"fmt"
	"log"
	"time"

	img "github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	bird *bird
	bg   *sdl.Texture
	pipe *pipe
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
	p, err := newPipe(r)
	if err != nil {
		return nil, err
	}

	return &scene{bg: bg, bird: b, pipe: p}, nil
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
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
			case <-tick:
				s.update()
				if s.bird.isDead() {
					drawTitle(r, "Game Over")
					time.Sleep(time.Second)
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

	case *sdl.MouseButtonEvent:
		s.bird.jump()
	case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.TouchFingerEvent:
	default:
		log.Printf("unknown event %T", event)
		return false
	}
	return false
}

func (s *scene) update() {
	s.bird.update()
	s.pipe.update()
	s.bird.touch(s.pipe)
}

func (s *scene) restart() {
	s.bird.restart()
	s.pipe.restart()

}
func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()
	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background:%v", err)
	}
	if err := s.bird.paint(r); err != nil {
		return err
	}
	if err := s.pipe.paint(r); err != nil {
		return err
	}
	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
	s.pipe.destroy()
}
