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

package internal

import (
	"github.com/finleap-connect/monoskope/internal/commandhandler"
	"github.com/finleap-connect/monoskope/internal/eventstore"
	"github.com/finleap-connect/monoskope/internal/gateway"
	"github.com/finleap-connect/monoskope/internal/queryhandler"
	"github.com/finleap-connect/monoskope/internal/test"
)

type TestEnv struct {
	*test.TestEnv
	gatewayTestEnv        *gateway.TestEnv
	eventStoreTestEnv     *eventstore.TestEnv
	queryHandlerTestEnv   *queryhandler.TestEnv
	commandHandlerTestEnv *commandhandler.TestEnv
}

func NewTestEnv(testEnv *test.TestEnv) (*TestEnv, error) {
	var err error
	env := &TestEnv{
		TestEnv: testEnv,
	}

	env.gatewayTestEnv, err = gateway.NewTestEnvWithParent(testEnv)
	if err != nil {
		return nil, err
	}

	env.eventStoreTestEnv, err = eventstore.NewTestEnvWithParent(testEnv)
	if err != nil {
		return nil, err
	}

	env.queryHandlerTestEnv, err = queryhandler.NewTestEnvWithParent(testEnv, env.eventStoreTestEnv, env.gatewayTestEnv)
	if err != nil {
		return nil, err
	}

	env.commandHandlerTestEnv, err = commandhandler.NewTestEnv(env.eventStoreTestEnv, env.gatewayTestEnv)
	if err != nil {
		return nil, err
	}

	return env, nil
}

func (env *TestEnv) Shutdown() error {
	if err := env.queryHandlerTestEnv.Shutdown(); err != nil {
		return err
	}

	if err := env.commandHandlerTestEnv.Shutdown(); err != nil {
		return err
	}

	if err := env.eventStoreTestEnv.Shutdown(); err != nil {
		return err
	}

	if err := env.gatewayTestEnv.Shutdown(); err != nil {
		return err
	}
	return nil
}
