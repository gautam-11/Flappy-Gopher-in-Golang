package main

import (
	"fmt"
	"sync"

	img "github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type pipe struct {
	mu       sync.RWMutex
	texture  *sdl.Texture
	x        int32
	h        int32
	w        int32
	speed    int32
	inverted bool
}

func newPipe(r *sdl.Renderer) (*pipe, error) {

	texture, err := img.LoadTexture(r, "res/imgs/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("could not load pipe image: %v", err)
	}

	return &pipe{
		texture:  texture,
		x:        400,
		h:        300,
		w:        50,
		speed:    1,
		inverted: false,
	}, nil
}
func (p *pipe) paint(r *sdl.Renderer) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	rect := &sdl.Rect{X: p.x, Y: (600 - p.h) - 43/2, W: p.w, H: p.h}
	if err := r.Copy(p.texture, nil, rect); err != nil {
		return fmt.Errorf("could not copy background:%v", err)
	}
	return nil
}

func (p *pipe) restart() {
	p.mu.RLock()
	defer p.mu.RUnlock()
	p.x = 400
}

func (p *pipe) update() {
	p.mu.RLock()
	defer p.mu.RUnlock()

	p.x -= p.speed
}

func (p *pipe) destroy() {
	p.mu.RLock()
	defer p.mu.RUnlock()
	p.texture.Destroy()
}
