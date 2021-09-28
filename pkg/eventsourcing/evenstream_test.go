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

package eventsourcing

import (
	"errors"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EventStream", func() {
	It("can stream events", func() {
		eventStream := NewEventStream()

		go func() {
			defer eventStream.Done()
			for i := 0; i < 3; i++ {
				eventStream.Send(event{})
			}
		}()

		received := 0
		for {
			event, err := eventStream.Receive()
			if err == io.EOF {
				break
			}
			Expect(err).To(Not(HaveOccurred()))
			Expect(event).To(Not(BeNil()))
			received++
		}
		Expect(received).To(Equal(3))
	})
	It("can handle errors", func() {
		eventStream := NewEventStream()

		go func() {
			defer eventStream.Done()
			eventStream.Error(errors.New("test"))
		}()

		for {
			event, err := eventStream.Receive()
			if err == io.EOF {
				break
			}
			Expect(err).To(HaveOccurred())
			Expect(event).To(BeNil())
		}
	})
})
