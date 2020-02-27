package elk

import "errors"

var ErrCircularDependency = errors.New("circular dependency")

var ErrTaskNotFound = errors.New("task not found")
