package strategy

import (
	"fmt"
)

type StrategyFactory func() Strategy

var allStrategies = map[string]StrategyFactory{
	"nop":    newNopStrategy,
	"greedy": newGreedyStrategy,
}

// Get a new instance of the Strategy with the given name.
func ForName(name string) (Strategy, error) {
	factory, ok := allStrategies[name]
	if !ok {
		return nil, fmt.Errorf("unknown strategy: %v", name)
	}

	s := factory()
	return s, nil
}
