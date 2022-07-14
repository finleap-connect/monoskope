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
	GatewayTestEnv        *gateway.TestEnv
	EventStoreTestEnv     *eventstore.TestEnv
	QueryHandlerTestEnv   *queryhandler.TestEnv
	CommandHandlerTestEnv *commandhandler.TestEnv
}

func NewTestEnv(testEnv *test.TestEnv) (env *TestEnv, err error) {
	env = &TestEnv{
		TestEnv: testEnv,
	}
	err = env.setup()
	return
}

func (env *TestEnv) setup() (err error) {
	env.EventStoreTestEnv, err = eventstore.NewTestEnvWithParent(env.TestEnv)
	if err != nil {
		return
	}

	env.GatewayTestEnv, err = gateway.NewTestEnvWithParent(env.TestEnv, env.EventStoreTestEnv, true)
	if err != nil {
		return
	}

	env.QueryHandlerTestEnv, err = queryhandler.NewTestEnvWithParent(env.TestEnv, env.EventStoreTestEnv, env.GatewayTestEnv)
	if err != nil {
		return
	}

	env.CommandHandlerTestEnv, err = commandhandler.NewTestEnv(env.EventStoreTestEnv, env.GatewayTestEnv)
	if err != nil {
		return
	}
	return
}

func (env *TestEnv) Shutdown() error {
	if err := env.QueryHandlerTestEnv.Shutdown(); err != nil {
		return err
	}

	if err := env.CommandHandlerTestEnv.Shutdown(); err != nil {
		return err
	}

	if err := env.EventStoreTestEnv.Shutdown(); err != nil {
		return err
	}

	if err := env.GatewayTestEnv.Shutdown(); err != nil {
		return err
	}
	return nil
}
