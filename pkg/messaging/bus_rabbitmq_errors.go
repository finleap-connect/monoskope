package messaging

import "errors"

// ErrConfigNameRequired is when the config doesn't include a name.
var ErrConfigNameRequired = errors.New("name must not be empty")
