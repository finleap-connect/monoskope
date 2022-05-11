// Copyright 2022 Monoskope Authors
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

	"github.com/google/uuid"
)

// Projector is the interface for projectors.
type Projector interface {
	// NewProjection creates a new Projection of the type the Projector projects.
	NewProjection(uuid.UUID) Projection

	// Project updates the state of the projection according to the given event.
	Project(context.Context, Event, Projection) (Projection, error)
}
