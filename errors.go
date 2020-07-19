package coffeemachine

import (
	"github.com/anshal21/coffee-machine/lib/errors"
)

// Set of error codes thrown by the expression evaluation
const (
	// ErrInvalidRuleSet represents some error in the provided ruleset
	ErrInvalidRuleSet errors.ErrCode = "ErrInvalidRuleSet"
)
