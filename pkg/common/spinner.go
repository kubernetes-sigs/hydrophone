/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"
)

var spinnerFrames = []string{
	"⠈⠁", "⠈⠑", "⠈⠱", "⠈⡱", "⢀⡱", "⢄⡱", "⢄⡱", "⢆⡱", "⢎⡱", "⢎⡰",
	"⢎⡠", "⢎⡀", "⢎⠁", "⠎⠁", "⠊⠁",
}

type Spinner struct {
	stop        chan struct{}
	stopped     chan struct{}
	mu          *sync.Mutex
	running     bool
	writer      io.Writer
	ticker      *time.Ticker
	prefix      string
	suffix      string
	frameFormat string
}

func NewSpinner(w io.Writer) *Spinner {
	frameFormat := "\x1b[?7l\r%s%s%s\x1b[?7h"
	if runtime.GOOS == "windows" {
		frameFormat = "\r%s%s%s"
	}
	return &Spinner{
		stop:        make(chan struct{}, 1),
		stopped:     make(chan struct{}),
		mu:          &sync.Mutex{},
		writer:      w,
		frameFormat: frameFormat,
	}
}

func (s *Spinner) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		return
	}
	s.running = true
	s.ticker = time.NewTicker(time.Millisecond * 100)
	go func() {
		for {
			for _, frame := range spinnerFrames {
				select {
				case <-s.stop:
					func() {
						s.mu.Lock()
						defer s.mu.Unlock()
						s.ticker.Stop()
						s.running = false
						s.stopped <- struct{}{}
					}()
					return
				case <-s.ticker.C:
					func() {
						s.mu.Lock()
						defer s.mu.Unlock()
						fmt.Fprintf(s.writer, s.frameFormat, s.prefix, frame, s.suffix)
					}()
				}
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.stop <- struct{}{}
	s.mu.Unlock()
	<-s.stopped
}
