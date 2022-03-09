package eventformatter

import (
	"context"
	"fmt"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"io"
	"strings"
)


type EventFormatter interface {
	GetFormattedDetails(context.Context, *esApi.Event) (string, error)
}

type BaseEventFormatter struct {
	EsClient esApi.EventStoreClient
}

// TODO: find a better place
// TODO: ticket: domain -> snapshots (same idea as e.g. domain -> projections)

func (f *BaseEventFormatter) GetSnapshot(ctx context.Context, projector es.Projector, eventFilter *esApi.EventFilter) (es.Projection, error) {
	projection := projector.NewProjection(uuid.New())
	aggregateEvents, err := f.EsClient.Retrieve(ctx, eventFilter)
	if err != nil {
		return nil, err
	}

	for {
		e, err := aggregateEvents.Recv()
		if err == io.EOF{
			break
		}
		if err != nil {
			return nil, err
		}

		event, err := es.NewEventFromProto(e)
		if err != nil {
			return nil, err
		}

		projection, err = projector.Project(ctx, event, projection)
		if err != nil {
			return nil, err
		}

	}

	return projection, nil
}

func (f *BaseEventFormatter) AppendUpdate(field string, update string, old string, strBuilder *strings.Builder) {
	if update != "" {
		strBuilder.WriteString(fmt.Sprintf("\n- “%s“ to “%s“", field, update))
		if old != "" {
			strBuilder.WriteString(fmt.Sprintf(" from “%s“", old))
		}
	}
}

// TODO: add to event data in es

func (f *BaseEventFormatter) ToPortoFromEventData(eventData []byte) (proto.Message, error) {
	porto := &anypb.Any{}
	if err := protojson.Unmarshal(eventData, porto); err != nil {
		return nil, err
	}
	ed, err := porto.UnmarshalNew()
	if err != nil {
		return nil, err
	}
	return ed, nil
}