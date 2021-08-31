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

package eventhandler

import (
	"context"

	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type loggingEventHandler struct {
	log logger.Logger
}

// NewLoggingEventHandler creates an EventHandler which automates storing Events in the EventStore when a Logging has emitted any.
func NewLoggingEventHandler() *loggingEventHandler {
	return &loggingEventHandler{
		log: logger.WithName("loggingEventHandler"),
	}
}

// HandleEvent implements the HandleEvent method of the es.EventHandler interface.
func (m *loggingEventHandler) HandleEvent(ctx context.Context, event es.Event) error {
	m.log.Info("Handling event.", "event", event.String())
	return nil
}
