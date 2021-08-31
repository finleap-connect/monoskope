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

import "github.com/google/uuid"

// Projection is the interface for projections.
type Projection interface {
	// ID returns the ID of the Projection.
	ID() uuid.UUID
	// Version returns the version of the aggregate this Projection is based upon.
	Version() uint64
	// IncrementVersion increments the Version of the Projection.
	IncrementVersion()
}
