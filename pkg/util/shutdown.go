// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

const (
	// Initially nothing is planned
	unplanned int32 = iota
	// Then only once a signal handler is registered
	registered
	// If a shutdown signal is received, the state is `planned`
	planned
)

type ShutdownWaitGroup struct {
	state int32 // Atomic variable defining the current state (see consts above)
	sync.WaitGroup
}

func NewShutdownWaitGroup() *ShutdownWaitGroup {
	return &ShutdownWaitGroup{}
}

func (s *ShutdownWaitGroup) IsExpected() bool {
	// If the `state` is not `planned`, we are not expecting a shutdown
	return atomic.LoadInt32(&s.state) == planned
}

func (s *ShutdownWaitGroup) Expect() {
	atomic.StoreInt32(&s.state, planned)
}

func (s *ShutdownWaitGroup) RegisterSignalHandler(shutdownCallback func()) {
	// Change our internal state to `registered`, if this is called twice it panics!
	swapped := atomic.CompareAndSwapInt32(&s.state, unplanned, registered)
	if !swapped {
		panic("shutdown signal handler registered twice!?")
	}
	s.Add(1) // Increment wg for the signal routine
	go func() {
		sigint := make(chan os.Signal, 1)
		// Interrupt signal sent from terminal or on sigterm
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		signal.Notify(sigint, syscall.SIGQUIT)
		<-sigint
		// We received a signal, so let's shutdown
		logger.WithName("shutdown-handler").Info("received shutdown signal")
		// Let's set the atomic properly to indicate planned shutdown behavior
		swapped := atomic.CompareAndSwapInt32(&s.state, registered, planned)
		if !swapped {
			panic("signal was received but atomic had unexpected value")
		}
		// Call the shutdown callback
		shutdownCallback()
		s.Done() // Routine done, let wg know
	}()
}

// Wait for internal `sync.WorkGroup` to complete and return `true` or `false`,
// if not shutdown successfully in timeout-limit.
func (s *ShutdownWaitGroup) WaitOrTimeout(timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		s.Wait()
	}()
	select {
	case <-c:
		return true
	case <-time.After(timeout):
		return false
	}
}
