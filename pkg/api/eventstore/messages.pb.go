// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.4
// source: api/eventstore/messages.proto

package eventstore

import (
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Event describes anything that has happened in the system.
// An event type name should be in past tense and contain the intent
// (TenantUpdated). The event should contain all the data needed when
// applying/handling it.
// The combination of aggregate_type, aggregate_id and version is
// unique.
type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Type of the event
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// Timestamp of when the event was created
	Timestamp *timestamp.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// ID of the aggregate that the event should be applied to (UUID 128-bit
	// number)
	AggregateId string `protobuf:"bytes,3,opt,name=aggregate_id,json=aggregateId,proto3" json:"aggregate_id,omitempty"`
	// Type of the aggregate that the event can be applied to
	AggregateType string `protobuf:"bytes,4,opt,name=aggregate_type,json=aggregateType,proto3" json:"aggregate_type,omitempty"`
	// Strict monotone counter, per aggregate/aggregate_id relation
	AggregateVersion *wrappers.UInt64Value `protobuf:"bytes,5,opt,name=aggregate_version,json=aggregateVersion,proto3" json:"aggregate_version,omitempty"`
	// Event type specific event data
	Data *any.Any `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_eventstore_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_api_eventstore_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_api_eventstore_messages_proto_rawDescGZIP(), []int{0}
}

func (x *Event) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Event) GetTimestamp() *timestamp.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Event) GetAggregateId() string {
	if x != nil {
		return x.AggregateId
	}
	return ""
}

func (x *Event) GetAggregateType() string {
	if x != nil {
		return x.AggregateType
	}
	return ""
}

func (x *Event) GetAggregateVersion() *wrappers.UInt64Value {
	if x != nil {
		return x.AggregateVersion
	}
	return nil
}

func (x *Event) GetData() *any.Any {
	if x != nil {
		return x.Data
	}
	return nil
}

// Request to get Events from to the store
type EventFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to ByAggregate:
	//	*EventFilter_AggregateId
	//	*EventFilter_AggregateType
	ByAggregate isEventFilter_ByAggregate `protobuf_oneof:"by_aggregate"`
	// Filter events with a version >= min_version
	MinVersion *wrappers.UInt64Value `protobuf:"bytes,3,opt,name=min_version,json=minVersion,proto3" json:"min_version,omitempty"`
	// Filter events with a version <= max_version
	MaxVersion *wrappers.UInt64Value `protobuf:"bytes,4,opt,name=max_version,json=maxVersion,proto3" json:"max_version,omitempty"`
	// Filter events with a timestamp >= min_timestamp
	MinTimestamp *timestamp.Timestamp `protobuf:"bytes,7,opt,name=min_timestamp,json=minTimestamp,proto3" json:"min_timestamp,omitempty"`
	// Filter events with a timestamp <= max_timestamp
	MaxTimestamp *timestamp.Timestamp `protobuf:"bytes,8,opt,name=max_timestamp,json=maxTimestamp,proto3" json:"max_timestamp,omitempty"`
}

func (x *EventFilter) Reset() {
	*x = EventFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_eventstore_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventFilter) ProtoMessage() {}

func (x *EventFilter) ProtoReflect() protoreflect.Message {
	mi := &file_api_eventstore_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventFilter.ProtoReflect.Descriptor instead.
func (*EventFilter) Descriptor() ([]byte, []int) {
	return file_api_eventstore_messages_proto_rawDescGZIP(), []int{1}
}

func (m *EventFilter) GetByAggregate() isEventFilter_ByAggregate {
	if m != nil {
		return m.ByAggregate
	}
	return nil
}

func (x *EventFilter) GetAggregateId() *wrappers.StringValue {
	if x, ok := x.GetByAggregate().(*EventFilter_AggregateId); ok {
		return x.AggregateId
	}
	return nil
}

func (x *EventFilter) GetAggregateType() *wrappers.StringValue {
	if x, ok := x.GetByAggregate().(*EventFilter_AggregateType); ok {
		return x.AggregateType
	}
	return nil
}

func (x *EventFilter) GetMinVersion() *wrappers.UInt64Value {
	if x != nil {
		return x.MinVersion
	}
	return nil
}

func (x *EventFilter) GetMaxVersion() *wrappers.UInt64Value {
	if x != nil {
		return x.MaxVersion
	}
	return nil
}

func (x *EventFilter) GetMinTimestamp() *timestamp.Timestamp {
	if x != nil {
		return x.MinTimestamp
	}
	return nil
}

func (x *EventFilter) GetMaxTimestamp() *timestamp.Timestamp {
	if x != nil {
		return x.MaxTimestamp
	}
	return nil
}

type isEventFilter_ByAggregate interface {
	isEventFilter_ByAggregate()
}

type EventFilter_AggregateId struct {
	// Filter events by aggregate_id
	AggregateId *wrappers.StringValue `protobuf:"bytes,1,opt,name=aggregate_id,json=aggregateId,proto3,oneof"`
}

type EventFilter_AggregateType struct {
	// Filter events for a specific aggregate type
	AggregateType *wrappers.StringValue `protobuf:"bytes,2,opt,name=aggregate_type,json=aggregateType,proto3,oneof"`
}

func (*EventFilter_AggregateId) isEventFilter_ByAggregate() {}

func (*EventFilter_AggregateType) isEventFilter_ByAggregate() {}

var File_api_eventstore_messages_proto protoreflect.FileDescriptor

var file_api_eventstore_messages_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x94, 0x02, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12,
	0x21, 0x0a, 0x0c, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65,
	0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x67, 0x67, 0x72,
	0x65, 0x67, 0x61, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x49, 0x0a, 0x11, 0x61, 0x67, 0x67,
	0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x10, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xa7,
	0x03, 0x0a, 0x0b, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x41,
	0x0a, 0x0c, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x48, 0x00, 0x52, 0x0b, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x49,
	0x64, 0x12, 0x45, 0x0a, 0x0e, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x48, 0x00, 0x52, 0x0d, 0x61, 0x67, 0x67, 0x72, 0x65,
	0x67, 0x61, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x3d, 0x0a, 0x0b, 0x6d, 0x69, 0x6e, 0x5f,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x55, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x6d, 0x69, 0x6e,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3d, 0x0a, 0x0b, 0x6d, 0x61, 0x78, 0x5f, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55,
	0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x6d, 0x61, 0x78, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3f, 0x0a, 0x0d, 0x6d, 0x69, 0x6e, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0c, 0x6d, 0x69, 0x6e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x3f, 0x0a, 0x0d, 0x6d, 0x61, 0x78, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0c, 0x6d, 0x61, 0x78, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x0e, 0x0a, 0x0c, 0x62, 0x79, 0x5f, 0x61,
	0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x42, 0x45, 0x5a, 0x43, 0x67, 0x69, 0x74, 0x6c,
	0x61, 0x62, 0x2e, 0x66, 0x69, 0x67, 0x6f, 0x2e, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x73, 0x2f,
	0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f,
	0x70, 0x65, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x73, 0x6b, 0x6f, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_eventstore_messages_proto_rawDescOnce sync.Once
	file_api_eventstore_messages_proto_rawDescData = file_api_eventstore_messages_proto_rawDesc
)

func file_api_eventstore_messages_proto_rawDescGZIP() []byte {
	file_api_eventstore_messages_proto_rawDescOnce.Do(func() {
		file_api_eventstore_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_eventstore_messages_proto_rawDescData)
	})
	return file_api_eventstore_messages_proto_rawDescData
}

var file_api_eventstore_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_eventstore_messages_proto_goTypes = []interface{}{
	(*Event)(nil),                // 0: eventstore.Event
	(*EventFilter)(nil),          // 1: eventstore.EventFilter
	(*timestamp.Timestamp)(nil),  // 2: google.protobuf.Timestamp
	(*wrappers.UInt64Value)(nil), // 3: google.protobuf.UInt64Value
	(*any.Any)(nil),              // 4: google.protobuf.Any
	(*wrappers.StringValue)(nil), // 5: google.protobuf.StringValue
}
var file_api_eventstore_messages_proto_depIdxs = []int32{
	2, // 0: eventstore.Event.timestamp:type_name -> google.protobuf.Timestamp
	3, // 1: eventstore.Event.aggregate_version:type_name -> google.protobuf.UInt64Value
	4, // 2: eventstore.Event.data:type_name -> google.protobuf.Any
	5, // 3: eventstore.EventFilter.aggregate_id:type_name -> google.protobuf.StringValue
	5, // 4: eventstore.EventFilter.aggregate_type:type_name -> google.protobuf.StringValue
	3, // 5: eventstore.EventFilter.min_version:type_name -> google.protobuf.UInt64Value
	3, // 6: eventstore.EventFilter.max_version:type_name -> google.protobuf.UInt64Value
	2, // 7: eventstore.EventFilter.min_timestamp:type_name -> google.protobuf.Timestamp
	2, // 8: eventstore.EventFilter.max_timestamp:type_name -> google.protobuf.Timestamp
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_api_eventstore_messages_proto_init() }
func file_api_eventstore_messages_proto_init() {
	if File_api_eventstore_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_eventstore_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_eventstore_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventFilter); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_api_eventstore_messages_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*EventFilter_AggregateId)(nil),
		(*EventFilter_AggregateType)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_eventstore_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_eventstore_messages_proto_goTypes,
		DependencyIndexes: file_api_eventstore_messages_proto_depIdxs,
		MessageInfos:      file_api_eventstore_messages_proto_msgTypes,
	}.Build()
	File_api_eventstore_messages_proto = out.File
	file_api_eventstore_messages_proto_rawDesc = nil
	file_api_eventstore_messages_proto_goTypes = nil
	file_api_eventstore_messages_proto_depIdxs = nil
}
