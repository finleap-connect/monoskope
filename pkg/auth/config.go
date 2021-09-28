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

package auth

type Config struct {
	// Hostname of the redirect URL.
	RedirectURLHostname string

	// Candidates of hostname and port which the local server binds to.
	// You can set port number to 0 to allocate a free port.
	// If multiple addresses are given, it will try the ports in order.
	// If nil or an empty slice is given, it defaults to "127.0.0.1:0" i.e. a free port.
	LocalServerBindAddress []string

	// Response HTML body on authorization completed.
	// Default to DefaultLocalServerSuccessHTML.
	LocalServerSuccessHTML string

	// A channel to send its URL when the local server is ready. Default to none.
	LocalServerReadyChan chan<- string
}
