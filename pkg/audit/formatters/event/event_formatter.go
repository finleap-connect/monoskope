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

package event

import (
	"context"
	"fmt"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/audit/formatters"
	"strings"
)

// EventFormatter is the interface definition for all event formatters
type EventFormatter interface {
	// GetFormattedDetails formats a given event in a human-readable format
	GetFormattedDetails(context.Context, *esApi.Event) (string, error)
}

// EventFormatterBase is the base implementation for all event formatters
type EventFormatterBase struct {
	*formatters.FormatterBase
}

// AppendUpdate appends updates to a string builder in human-readable format
func (f *EventFormatterBase) AppendUpdate(field string, update string, old string, strBuilder *strings.Builder) {
	if update != "" {
		strBuilder.WriteString(fmt.Sprintf("\n- “%s“ to “%s“", field, update))
		if old != "" {
			strBuilder.WriteString(fmt.Sprintf(" from “%s“", old))
		}
	}
}
