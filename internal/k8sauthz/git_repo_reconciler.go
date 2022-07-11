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

package k8sauthz

// GitRepoReconciler reconciles the resources within the target repo to match the expected state.
type GitRepoReconciler struct {
}

// NewGitRepoReconciler creates a new GitRepoReconciler configured via the given config.
func NewGitRepoReconciler(config *GitRepoReconcilerConfig) (*GitRepoReconciler, error) {
	panic("not implemented")
}
