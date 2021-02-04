package event_sourcing

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventdata/test"
)

var _ = Describe("handler", func() {
	eventCounter := 0

	createTestEventData := func(something string) EventData {
		ed, err := ToEventDataFromProto(&test.TestEventData{Hello: something})
		Expect(err).ToNot(HaveOccurred())
		return ed
	}

	createEvent := func() Event {
		data := createTestEventData("world!")
		event := NewEvent(testEventType, data, time.Now().UTC(), testAggregateType, uuid.New(), uint64(eventCounter))
		eventCounter++
		return event
	}

	It("can chain command handler", func() {
		cmd := &testCommand{
			TestCommandData: commands.TestCommandData{
				Test:      "test",
				TestCount: 0,
			},
		}
		handlerChain := ChainCommandHandler(
			&testCommandHandler{val: 1},
			&testCommandHandler{val: 2},
			&testCommandHandler{val: 3},
		)
		err := handlerChain.HandleCommand(context.Background(), cmd)
		Expect(err).ToNot(HaveOccurred())
		Expect(cmd.TestCommandData.TestCount).To(BeNumerically("==", 3))
		Expect(cmd.TestCommandData.Test).To(Equal("test123"))
	})
	It("can chain event handler", func() {
		event := createEvent()
		handlerChain := ChainEventHandler(
			&testEventHandler{},
			&testEventHandler{},
			&testEventHandler{},
		)
		err := handlerChain.HandleEvent(context.Background(), event)
		Expect(err).ToNot(HaveOccurred())
	})
})
