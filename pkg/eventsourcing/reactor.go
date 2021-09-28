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
	"context"
)

// Reactor is the interface for reactors.
type Reactor interface {
	// HandleEvent handles a given event send 0..* Events through the given channel in reaction or an error.
	// Attention: The reactor is responsible for closing the channel if no further events will be send to that channel.
	HandleEvent(ctx context.Context, event Event, events chan<- Event) error
}
