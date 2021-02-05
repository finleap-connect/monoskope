package event_sourcing

// Projection is the interface for projections.
type Projection interface {
	// ID returns the ID of the Projection.
	GetId() string
	// AggregateVersion is the version of the Aggregate this Projection is based upon.
	GetAggregateVersion() uint64
}
