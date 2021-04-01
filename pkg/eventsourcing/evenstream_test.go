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
