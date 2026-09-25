package main

import (
	"errors"
	"io"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type Player struct {
	mu     sync.Mutex
	cmd    *exec.Cmd
	done   chan struct{}
	url    string
	paused bool
}

func (p *Player) Play(url string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, err := exec.LookPath("mpv"); err != nil {
		return errors.New("mpv not found. install it first (brew install mpv)")
	}
	p.kill()

	cmd := exec.Command("mpv", "--no-video", "--no-terminal", "--really-quiet", url)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Start(); err != nil {
		return err
	}

	done := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(done)
	}()

	p.cmd = cmd
	p.done = done
	p.url = url
	p.paused = false
	return nil
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.kill()
	p.paused = false
	p.url = ""
}

func (p *Player) Toggle() error {
	p.mu.Lock()
	paused := p.paused
	url := p.url
	p.mu.Unlock()

	if paused {
		if url == "" {
			return nil
		}
		return p.Play(url)
	}
	p.Pause()
	return nil
}

func (p *Player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.cmd == nil {
		return
	}
	p.kill()
	p.paused = true
}

func (p *Player) Paused() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.paused
}

func (p *Player) kill() {
	if p.cmd == nil || p.cmd.Process == nil {
		p.cmd = nil
		p.done = nil
		return
	}
	_ = p.cmd.Process.Signal(syscall.SIGTERM)
	select {
	case <-p.done:
	case <-time.After(3 * time.Second):
		_ = p.cmd.Process.Kill()
		<-p.done
	}
	p.cmd = nil
	p.done = nil
}
